package expression

import (
	"ExprCalc/internal/models"
	"ExprCalc/pkg/broker"
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/twharmon/gouid"
)

type worker struct {
	id        string
	lastTouch time.Time
	handler   func(*models.Expression)
	input     <-chan amqp.Delivery
	rabbit    *broker.RabbitMQ
	ctx       context.Context
	cancel    context.CancelFunc
}

func newWorker(rabbit *broker.RabbitMQ, handler func(*models.Expression), input <-chan amqp.Delivery) *worker {
	ctx, cancel := context.WithCancel(context.Background())
	id := gouid.Bytes(32)
	worker := &worker{
		rabbit:  rabbit,
		id:      id.String(),
		ctx:     ctx,
		cancel:  cancel,
		handler: handler,
	}

	worker.input = input
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
		case input := <-w.input:
			expr := new(models.Expression)
			err := expr.UnmarshalBinary(input.Body)
			expr.Err = err

			w.handler(expr)

			body, err := expr.MarshalBinary()
			fmt.Println(string(body), "должно выйти")

			if err != nil {
				expr.Err = err
			}
			err = w.rabbit.Ch.PublishWithContext(w.ctx, "", input.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "application/json",
				Body:          body,
				CorrelationId: input.CorrelationId,
			})
			if err != nil {

			}
			w.lastTouch = time.Now()
			input.Ack(false)
		case <-w.ctx.Done():
			return
		}

	}
}
