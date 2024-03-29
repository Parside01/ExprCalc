package expression

import (
	"ExprCalc/internal/models"
	"ExprCalc/pkg/broker"
	"ExprCalc/pkg/config"
	"ExprCalc/pkg/repository/redisdb"
	"context"
	"fmt"
	"time"

	"github.com/Knetic/govaluate"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var (
	ErrInvalidExpression = fmt.Errorf("invalid expression")
)

type ExpressionService struct {
	logger     *zap.Logger
	config     *config.ExpressionServiceConfig
	rabbit     *broker.RabbitMQ
	workers    map[string]*worker
	cache      *redisdb.RedisDB
	listenExpr <-chan amqp.Delivery
}

/*
*	Возможно надо было сделать в отдельном модуле.
*	Всем же так нравится горизонтальное масштабирование.
*	Но я пожалею вашу оперативку и оставлю так (1 контейнер будет +- 600мб)
 */
func NewExpressionService(logger *zap.Logger, cfg *config.ExpressionServiceConfig, rabbit *broker.RabbitMQ, cache *redisdb.RedisDB) (*ExpressionService, error) {
	service := &ExpressionService{
		logger:  logger,
		workers: make(map[string]*worker),
		config:  cfg,
		cache:   cache,
		rabbit:  rabbit,
	}

	// Что-то пошло не так == отвал башки.
	err := service.start()
	if err != nil {
		logger.Fatal("ExpressionService.NewExpressionService: failed to start expression service", zap.Error(err))
	}
	return service, nil
}

func (e *ExpressionService) start() error {
	err := e.setupRabbit()
	e.setupWorkers()

	return err
}

/*
*	Запускаем воркеры, которые представляют оболочку надо горутинами, которые обрабатывают то что приходит в RabbitMQ.
 */
func (e *ExpressionService) setupWorkers() {
	for i := 0; i < e.config.GourutinesCount; i++ {
		worker := newWorker(e.logger, e.rabbit, e.handle, e.listenExpr, time.Duration(e.config.WorkerInfoUpdate*int(time.Second)), e.cache)
		if worker == nil {
			i--
			continue
		}
		e.workers[worker.id] = worker
	}
}

/*
*	Вы можете спросить меня нафига это надо везде.
*	Я не придумал как это сделать в виде короткой функции внутри broker.
*	По факту просто настраиваем соединение с RabbitMQ.
*	Если что-то поменять то оно отвалится. Сколько я раз пожалел что не gRPC
 */
func (e *ExpressionService) setupRabbit() error {
	q, err := e.rabbit.Ch.QueueDeclare(e.config.ExpressionQueue, false, false, false, false, nil)
	if err != nil {
		return err
	}

	err = e.rabbit.Ch.QueueBind(q.Name, e.config.RouteKey, e.config.Exchange, false, nil)
	if err != nil {
		return err
	}

	ch, err := e.rabbit.Ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	e.listenExpr = ch
	return nil
}

func (e *ExpressionService) Stop() {
	for _, worker := range e.workers {
		worker.cancel()
	}
}

func (e *ExpressionService) handle(expr *models.Expression) {
	start := time.Now()
	re, err := govaluate.NewEvaluableExpression(expr.Expression)
	if err != nil {
		expr.Err = err
		return
	}
	result, err := re.Evaluate(nil)
	if err != nil {
		expr.Err = err
		return
	}

	/*
	*	Хах. Мы уже все подсчитали и отдали, но запишим это в базу только через некоторое время.
	*	Только после этого сервак будет знать что все готово).
	 */
	go func(ctx context.Context, start time.Time, expr *models.Expression, res float64) {
		select {
		case <-time.After(time.Duration(expr.ExpectExucuteTime) * time.Millisecond):
			expr.ExecuteTime = time.Since(start).Milliseconds()
			expr.Result = res
			expr.IsDone = true
			e.cache.WriteCache(ctx, expr.Expression, expr)
		case <-ctx.Done():
			return
		}
	}(context.Background(), start, expr, result.(float64))
}
