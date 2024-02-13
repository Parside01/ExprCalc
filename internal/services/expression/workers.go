package expression

import (
	"ExprCalc/internal/models"
	"ExprCalc/pkg/broker"
	"ExprCalc/pkg/repository/redisdb"
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/twharmon/gouid"
	"go.uber.org/zap"
)

type worker struct {
	logger    *zap.Logger
	id        string
	lastTouch time.Time
	handler   func(*models.Expression)
	cache     *redisdb.RedisDB
	inputExpr <-chan amqp.Delivery
	rabbit    *broker.RabbitMQ
	ctx       context.Context    //да да не очень хорошая практика.
	cancel    context.CancelFunc // и это тоже.
	writeInfo <-chan amqp.Delivery
}

func newWorker(logger *zap.Logger, rabbit *broker.RabbitMQ, handler func(*models.Expression), input <-chan amqp.Delivery) *worker {
	ctx, cancel := context.WithCancel(context.Background())
	id := gouid.Bytes(16)
	worker := &worker{
		logger:  logger,
		rabbit:  rabbit,
		id:      id.String(),
		ctx:     ctx,
		cancel:  cancel,
		handler: handler,
	}

	worker.inputExpr = input
	go worker.startLoop()
	return worker
}

/*
*	Логика не очень сложная.
* 	Просто запускаем цикл обработки сообщений в горутинах.
*	После добработки отправляем тому кто прислал через @RoutingKey и @CorrelationId.
 */
func (w *worker) startLoop() {
	for {
		select {
		case input := <-w.inputExpr:
			w.proccessExpression(input)
		case <-w.ctx.Done():
			return
		}

	}
}

func (w *worker) proccessExpression(input amqp.Delivery) {
	w.logger.Info(string(input.Body), zap.String(fmt.Sprintf("Пришло в worker"), w.id))

	expr := new(models.Expression)
	err := expr.UnmarshalBinary(input.Body)
	expr.Err = err
	expr.WorkerID = w.id

	w.handler(expr)

	body, err := expr.MarshalBinary()
	if err != nil {
		expr.Err = err
	}
	err = w.rabbit.Ch.PublishWithContext(w.ctx, "", input.ReplyTo, false, false, amqp.Publishing{
		ContentType:   "application/json",
		Body:          body,
		CorrelationId: input.CorrelationId,
	})
	if err != nil {
		w.logger.Error("ExpressionController.calcHandler: failed to publish message", zap.Error(err))
	}
	w.lastTouch = time.Now()

	input.Ack(false)
}

func (w *worker) processExprInRedis(expr *models.Expression) (*models.Expression, error) {
	ok, err := w.cache.IsExist(w.ctx, expr.Expression)
	if err != nil {
		return nil, err
	}

	if ok == 1 {
		res, err := w.cache.GetCache(w.ctx, expr.Expression)
		if err != nil {
			w.logger.Error("ExpressionService.processExprInRedis: failed to get existing cache", zap.Error(err))
			return nil, err
		}

		err = json.Unmarshal([]byte(res.(string)), &expr)
		if err != nil {
			w.logger.Error("ExpressionService.processExprInRedis: failed to unmarshal expression", zap.Error(err))
			return nil, err
		}
		return expr, nil
	}
	return nil, nil
}
