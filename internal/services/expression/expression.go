package expression

import (
	"ExprCalc/internal/models"
	"ExprCalc/pkg/broker"
	"ExprCalc/pkg/config"
	"fmt"
	"time"

	"github.com/Maldris/mathparse"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var (
	ErrInvalidExpression = fmt.Errorf("invalid expression")
)

type ExpressionService struct {
	logger  *zap.Logger
	config  *config.ExpressionServiceConfig
	rabbit  *broker.RabbitMQ
	workers map[string]*worker
	listen  <-chan amqp.Delivery
}

func NewExpressionService(logger *zap.Logger, cfg *config.ExpressionServiceConfig, rabbit *broker.RabbitMQ) (*ExpressionService, error) {
	service := &ExpressionService{
		logger:  logger,
		workers: make(map[string]*worker),
		config:  cfg,
		rabbit:  rabbit,
	}
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
func (e *ExpressionService) setupWorkers() {
	for i := 0; i < e.config.GourutinesCount; i++ {
		worker := newWorker(e.rabbit, e.parse, e.listen)
		if worker == nil {
			i--
			continue
		}
		e.workers[worker.id] = worker
	}
}

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

	e.listen = ch
	return nil
}

func (e *ExpressionService) Stop() {
	for _, worker := range e.workers {
		worker.cancel()
	}
}

func (e *ExpressionService) parse(expr *models.Expression) {
	start := time.Now()
	parser := mathparse.NewParser(expr.Expression)
	parser.Resolve()
	if !parser.FoundResult() {
		expr.Err = ErrInvalidExpression
	}
	expr.ExecuteTime = time.Since(start)
	res := parser.GetValueResult()
	expr.Result = res
}
