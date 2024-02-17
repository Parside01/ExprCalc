package broker

import (
	"ExprCalc/pkg/config"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	Conn   *amqp.Connection
	Ch     *amqp.Channel
	config *config.RabbitMQConfig
	logger *zap.Logger
}

func NewRabbit(logger *zap.Logger, config *config.RabbitMQConfig) *RabbitMQ {
	rabbit := &RabbitMQ{
		logger: logger,
		config: config,
	}

	rabbit.start()
	rabbit.logger.Info("RabbitMQ connected to", zap.String("uri", config.URI))
	return rabbit
}

func (r *RabbitMQ) start() {
	conn, err := amqp.Dial(r.config.URI)
	if err != nil {
		r.logger.Error("failed to connect to rabbit", zap.Error(err))
	}
	r.Conn = conn
	ch, err := r.Conn.Channel()
	if err != nil {
		r.logger.Error("failed to open a channel", zap.Error(err))
	}
	r.Ch = ch
}
