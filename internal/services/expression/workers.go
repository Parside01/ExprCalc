package expression

import (
	"ExprCalc/internal/models"
	"ExprCalc/pkg/broker"
	"ExprCalc/pkg/repository/redisdb"
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sourcegraph/conc"
	"github.com/twharmon/gouid"
	"go.uber.org/zap"
)

type worker struct {
	id         string
	isEmploy   bool
	lastTouch  time.Time
	currJob    string
	logger     *zap.Logger
	wg         *conc.WaitGroup //так на будущее
	handler    func(*models.Expression)
	cache      *redisdb.RedisDB
	inputExpr  <-chan amqp.Delivery
	rabbit     *broker.RabbitMQ
	infoUpdate time.Duration
	ctx        context.Context    //да да не очень хорошая практика.
	cancel     context.CancelFunc // и это тоже.
}

func newWorker(logger *zap.Logger, rabbit *broker.RabbitMQ, handler func(*models.Expression), input <-chan amqp.Delivery, infoUpdate time.Duration, cache *redisdb.RedisDB) *worker {
	ctx, cancel := context.WithCancel(context.Background())
	id := gouid.Bytes(16)
	worker := &worker{
		logger:     logger,
		rabbit:     rabbit,
		id:         id.String(),
		cache:      cache,
		ctx:        ctx,
		cancel:     cancel,
		handler:    handler,
		wg:         conc.NewWaitGroup(),
		infoUpdate: infoUpdate,
		isEmploy:   false,
		currJob:    "",
	}

	worker.inputExpr = input
	go worker.startExprLoop()
	go worker.startCacheLoop(time.NewTicker(infoUpdate))

	err := worker.cache.WriteCacheWithTTL(worker.ctx, fmt.Sprintf("woker-%s", worker.id), &workerInfo{
		WorkerID:   worker.id,
		LastTouch:  time.Now().String(),
		IsEmploy:   worker.isEmploy,
		CurrentJob: worker.currJob,
	}, infoUpdate)
	fmt.Println(infoUpdate)
	if err != nil {
		worker.logger.Error(fmt.Sprintf("worker.startCacheLoop: failed to write cache in worker %s", worker.id), zap.Error(err))
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
			err := w.cache.WriteCacheWithTTL(w.ctx, fmt.Sprintf("woker-%s", w.id), &workerInfo{
				WorkerID:   w.id,
				LastTouch:  w.lastTouch.String(),
				IsEmploy:   w.isEmploy,
				CurrentJob: w.currJob,
			}, w.infoUpdate)
			if err != nil {
				w.logger.Error(fmt.Sprintf("worker.startCacheLoop: failed to write cache in worker %s", w.id), zap.Error(err))
			}
		case <-w.ctx.Done():
			return
		}
	}
}

func (w *worker) proccessExpression(input amqp.Delivery) {
	w.logger.Info(string(input.Body), zap.String("Пришло в worker", w.id))

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
