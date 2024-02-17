package controllers

import (
	"ExprCalc/internal/models"
	"ExprCalc/internal/server/middleware"
	"ExprCalc/pkg/broker"
	"ExprCalc/pkg/config"
	"ExprCalc/pkg/repository/redisdb"
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/twharmon/gouid"
	"go.uber.org/zap"
)

type ExpressionController struct {
	logger     *zap.Logger
	config     *config.ExpressionServiceConfig
	rabbit     *broker.RabbitMQ
	cache      *redisdb.RedisDB
	resultExpr <-chan amqp.Delivery
}

type Request struct {
	Expression       string `json:"expr"`
	MultiplyTime     int64  `json:"*"`
	AddTime          int64  `json:"+"`
	DivideTime       int64  `json:"/"`
	SubtractTime     int64  `json:"-"`
	PowTime          int64  `json:"**"`
	DivRemainderTime int64  `json:"%"`
}

type Response struct {
	Expr *models.Expression `json:"expr"`
	Err  error              `json:"err"`
	Ok   bool               `json:"ok"`
}

func NewExpressionController(logger *zap.Logger, config *config.ExpressionServiceConfig, rabbit *broker.RabbitMQ, cache *redisdb.RedisDB) *ExpressionController {
	c := &ExpressionController{
		logger: logger,
		config: config,
		rabbit: rabbit,
		cache:  cache,
	}
	err := c.setupRabbit()
	if err != nil {
		logger.Fatal("ExpressionController.NewExpressionController: failed to start expression service", zap.Error(err))
	}

	return c
}

func (e *ExpressionController) GetGroup() string {
	return "/expr"
}
func (e *ExpressionController) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.CorseDisable(),
	}
}

func (e *ExpressionController) GetHandlers() []ControllerHandler {
	return []ControllerHandler{
		&Handler{
			Method:  "POST",
			Path:    "/calc",
			Handler: e.calcHandler,
		},
		&Handler{
			Method:  "GET",
			Path:    "/getWorkersInfo",
			Handler: e.getWorkersInfo,
		},
		&Handler{
			Method:  "GET",
			Path:    "/getAllExpressions",
			Handler: e.getAllExpressions,
		},
	}
}

/*
*	Принимаем сообщение. Нужно дописать его валидацию, чтобы отдавать нужные ошибки.
*	С целью скорости отдачи ошибок при неверно введенных данных будем валидировать все на фронте.
*	Если оно есть в кэше то отдаем то что в кэше.
* 	Отпраяляем в обработчики. Получаем результат через засетапленный канал.
* 	Фиксируем результат.
 */
func (e *ExpressionController) calcHandler(c echo.Context) error {
	var req Request
	err := c.Bind(&req)
	if err != nil {
		e.logger.Error("ExpressionController.calcHandler: failed to bind request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, &Response{Err: err, Ok: false})
	}

	ok, err := e.cache.IsExist(c.Request().Context(), req.Expression)
	if err != nil {
		e.logger.Error("ExpressionController.calcHandler: failed to check cache", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, &Response{Err: err, Ok: false})
	}

	if ok == 1 {
		res, err := e.cache.GetCache(c.Request().Context(), req.Expression)
		if err != nil {
			e.logger.Error("ExpressionController.calcHandler: failed to get cache", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, &Response{Err: err, Ok: false})
		}

		var expr *models.Expression
		err = json.Unmarshal([]byte(res.(string)), &expr)
		if err != nil {
			e.logger.Error("ExpressionController.calcHandler: failed to unmarshal expression", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, &Response{Err: err, Ok: false})
		}

		return c.JSON(http.StatusOK, &Response{Expr: expr, Err: nil, Ok: true})
	}

	expr := models.NewExpression(req.Expression)
	expr.ExpectExucuteTime = CalcExecurteTime(req)

	body, err := expr.MarshalBinary()
	if err != nil {
		e.logger.Error("ExpressionController.calcHandler: failed to marshal expression", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, &Response{Err: err, Ok: false})
	}

	corrId := gouid.Bytes(32)

	err = e.rabbit.Ch.PublishWithContext(c.Request().Context(), e.config.Exchange, e.config.RouteKey, false, false, amqp.Publishing{
		ReplyTo:       e.config.ResultQueue,
		CorrelationId: corrId.String(),
		Body:          body,
		ContentType:   "application/json",
	})
	if err != nil {
		e.logger.Error("ExpressionController.calcHandler: failed to publish message", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, &Response{Err: err, Ok: false})
	}

	for {
		select {
		case msg := <-e.resultExpr:
			if msg.CorrelationId == corrId.String() {
				e.cache.WriteCache(c.Request().Context(), req.Expression, msg.Body)

				var expr *models.Expression

				err = json.Unmarshal(msg.Body, &expr)
				if err != nil {
					e.logger.Error("ExpressionController.calcHandler: failed to unmarshal expression", zap.Error(err))
					return c.JSON(http.StatusInternalServerError, &Response{Err: err, Ok: false})
				}

				expr.IsDone = true
				msg.Ack(false)

				return c.JSON(http.StatusOK, &Response{Expr: expr, Err: nil, Ok: true})
			}
		case <-c.Request().Context().Done():
			return c.Request().Context().Err()
		}
	}
}

/*
*	В каждом воркере есть штука которая записывает информацию о нем в редиску.
* 	Вот эту информацию находим и отдаем.
 */
func (e *ExpressionController) getWorkersInfo(c echo.Context) error {
	keys, err := e.cache.GetAllKeysByPattern(c.Request().Context(), "worker")
	if err != nil {
		e.logger.Error("ExpressionController.getWorkersInfo: failed to get workers info", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, &Response{Err: err, Ok: false})
	}

	res := []*models.WorkerInfo{}
	for _, i := range keys {
		var info *models.WorkerInfo
		obj, err := e.cache.GetCache(c.Request().Context(), i)
		if err != nil {
			e.logger.Error("ExpressionController.getWorkersInfo: failed to get workers info", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, &Response{Err: err, Ok: false})
		}

		err = json.Unmarshal([]byte(obj.(string)), &info)
		if err != nil {
			e.logger.Error("ExpressionController.getWorkersInfo: failed to unmarshal workers info", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, &Response{Err: err, Ok: false})
		}

		res = append(res, info)
	}
	return c.JSON(http.StatusOK, res)
}

func (e *ExpressionController) getAllExpressions(c echo.Context) error {
	take, err := strconv.Atoi(c.QueryParams().Get("take"))
	if err != nil {
		e.logger.Error("ExpressionController.getAllExpressions: failed to get take param", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, &Response{Err: err, Ok: false})
	}

	skip, err := strconv.Atoi(c.QueryParams().Get("skip"))
	if err != nil {
		e.logger.Error("ExpressionController.getAllExpressions: failed to get skip param", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, &Response{Err: err, Ok: false})
	}

	keys, err := e.cache.GetAllKeysByPattern(c.Request().Context(), "*")
	if err != nil {
		e.logger.Error("ExpressionController.getAllExpressions: failed to get all expressions", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, &Response{Err: err, Ok: false})
	}

	res := []*models.Expression{}
	s, t := 0, 0
	for _, i := range keys {
		if strings.Contains(i, "worker") {
			continue
		}
		if s <= skip {
			s++
			continue
		}
		t++

		f, err := e.cache.GetCache(c.Request().Context(), i)
		if err != nil {
			e.logger.Error("ExpressionController.getAllExpressions: failed to get expression", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, &Response{Err: err, Ok: false})
		}

		var expr *models.Expression
		err = json.Unmarshal([]byte(f.(string)), &expr)
		if err != nil {
			e.logger.Error("ExpressionController.getAllExpressions: failed to unmarshal expression", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, &Response{Err: err, Ok: false})
		}
		res = append(res, expr)

		if t >= take {
			break
		}
	}
	return c.JSON(http.StatusOK, res)
}

/*
*	Сетапим подключение к кролику для контроллер, но думаю что стоит как то передлать этот момент для всего в общем виде,
*	Ибо почти такой же код и в service, но без регистрации exchange.
* 	В очередь котрую тут добаляем будут приходить решенные выражения.
 */
func (e *ExpressionController) setupRabbit() error {

	err := e.rabbit.Ch.ExchangeDeclare(e.config.Exchange, "direct", true, false, false, false, nil)
	if err != nil {
		e.logger.Fatal("main.failed to declare exchange", zap.Error(err))
	}

	q, err := e.rabbit.Ch.QueueDeclare(e.config.ResultQueue, false, false, false, false, nil)
	if err != nil {
		e.logger.Fatal("main.failed to declare queue", zap.Error(err))
	}

	err = e.rabbit.Ch.QueueBind(q.Name, e.config.RouteKey, e.config.Exchange, false, nil)
	if err != nil {
		e.logger.Fatal("main.failed to bind queue", zap.Error(err))
	}

	out, err := e.rabbit.Ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		e.logger.Fatal("main.failed to consume queue", zap.Error(err))
	}

	e.resultExpr = out

	/*
	*	Вы думаете я знаю почему это так должно быть?
	* 	Нет я не знаю, но так оно заработало.
	 */
	go func() {
		for msg := range out {
			msg.Ack(false)
		}
	}()

	return nil
}

func CalcExecurteTime(req Request) int64 {
	arr := []int64{req.MultiplyTime, req.AddTime, req.DivideTime, req.SubtractTime, req.PowTime, req.DivRemainderTime}
	max := slices.Max(arr)

	return max
}
