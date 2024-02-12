package controllers

import (
	"ExprCalc/internal/models"
	"ExprCalc/pkg/broker"
	"ExprCalc/pkg/config"
	"ExprCalc/pkg/redisdb"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/twharmon/gouid"
	"go.uber.org/zap"
)

type ExpressionController struct {
	logger *zap.Logger
	config *config.ExpressionServiceConfig
	rabbit *broker.RabbitMQ
	cache  *redisdb.RedisDB
	listen <-chan amqp.Delivery
}

type Request struct {
	Expression string `json:"expr"`
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
	return nil
}

func (e *ExpressionController) GetHandlers() []ControllerHandler {
	return []ControllerHandler{
		&Handler{
			Method:  "POST",
			Path:    "/calc",
			Handler: e.calcHandler,
		},
	}
}

/*
*	Принимаем сообщение. Нужно дописать его валидацию, чтобы отдавать нужные ошибки.
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
		err = expr.UnmarshalBinary([]byte(res.(string)))
		if err != nil {
			e.logger.Error("ExpressionController.calcHandler: failed to unmarshal expression", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, &Response{Err: err, Ok: false})
		}

		return c.JSON(http.StatusOK, &Response{Expr: expr, Err: nil, Ok: true})
	}

	body, err := models.NewExpression(req.Expression).MarshalBinary()
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
		case msg := <-e.listen:
			if msg.CorrelationId == corrId.String() {
				e.cache.WriteCache(c.Request().Context(), req.Expression, msg.Body)

				var expr *models.Expression

				err = json.Unmarshal(msg.Body, &expr)
				if err != nil {
					e.logger.Error("ExpressionController.calcHandler: failed to unmarshal expression", zap.Error(err))
					return c.JSON(http.StatusInternalServerError, &Response{Err: err, Ok: false})
				}

				msg.Ack(false)
				return c.JSON(http.StatusOK, &Response{Expr: expr, Err: nil, Ok: true})
			}
		case <-c.Request().Context().Done():
			return c.Request().Context().Err()
		}
	}
}

/*
*	Сетапим подключение к кролику для контроллер, но думаю что стоит как то передлать этот момент для всего в общем виде,
*	Ибо почти такой же код и в service, но без регистрации exchange.
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

	e.listen = out

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
