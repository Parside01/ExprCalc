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
	"github.com/sourcegraph/conc"
	"github.com/twharmon/gouid"
	"go.uber.org/zap"
)

type worker struct {
	id        string
	isEmploy  bool
	lastTouch time.Time
	currJob   string
	logger    *zap.Logger
	wg        *conc.WaitGroup //так на будущее
	handler   func(*models.Expression)
	cache     *redisdb.RedisDB
	inputExpr <-chan amqp.Delivery
	rabbit    *broker.RabbitMQ
	ctx       context.Context    //да да не очень хорошая практика.
	cancel    context.CancelFunc // и это тоже.
}

func newWorker(logger *zap.Logger, rabbit *broker.RabbitMQ, handler func(*models.Expression), input <-chan amqp.Delivery, infoUpdate time.Duration, cache *redisdb.RedisDB) *worker {
	ctx, cancel := context.WithCancel(context.Background())
	id := gouid.Bytes(16)
	worker := &worker{
		logger:   logger,
		rabbit:   rabbit,
		id:       id.String(),
		cache:    cache,
		ctx:      ctx,
		cancel:   cancel,
		handler:  handler,
		wg:       conc.NewWaitGroup(),
		isEmploy: false,
		currJob:  "",
	}

	worker.inputExpr = input
	go worker.startExprLoop()
	go worker.startCacheLoop(time.NewTicker(infoUpdate))

	err := worker.cache.WriteCache(worker.ctx, fmt.Sprintf("woker-%s", worker.id), &workerInfo{
		WorkerID:   worker.id,
		LastTouch:  time.Now().String(),
		IsEmploy:   worker.isEmploy,
		CurrentJob: worker.currJob,
	})
	if err != nil {
		worker.logger.Error(fmt.Sprintf("worker.startCacheLoop: failed to write cache in worker", worker.id), zap.Error(err))
	}

	return worker
}

/*
*	Логика не очень сложная.
* 	Просто запускаем цикл обработки сообщений в горутинах.
*	После добработки отправляем тому кто прислал через @RoutingKey и @CorrelationId.
 */
func (w *worker) startExprLoop() {
	for {
		select {
		case input := <-w.inputExpr:
			w.proccessExpression(input)
		case <-w.ctx.Done():
			return
		}

	}
}

func (w *worker) startCacheLoop(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			err := w.cache.WriteCache(w.ctx, fmt.Sprintf("woker-%s", w.id), &workerInfo{
				WorkerID:   w.id,
				LastTouch:  w.lastTouch.String(),
				IsEmploy:   w.isEmploy,
				CurrentJob: w.currJob,
			})
			if err != nil {
				w.logger.Error(fmt.Sprintf("worker.startCacheLoop: failed to write cache in worker", w.id), zap.Error(err))
			}
		case <-w.ctx.Done():
			return
		}
	}
}

func (w *worker) proccessExpression(input amqp.Delivery) {
	w.logger.Info(string(input.Body), zap.String(fmt.Sprintf("Пришло в worker"), w.id))

	w.isEmploy = true
	defer w.onWaitState()

	expr := new(models.Expression)
	err := expr.UnmarshalBinary(input.Body)
	expr.Err = err

	w.currJob = expr.Expression

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
		return
	}

	input.Ack(false)
}

func (w *worker) onWaitState() {
	w.isEmploy = false
	w.lastTouch = time.Now()
	w.currJob = ""
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
