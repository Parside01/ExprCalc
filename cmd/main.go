package main

import (
	"ExprCalc/internal/server"
	"ExprCalc/internal/server/controllers"
	"ExprCalc/internal/services/expression"
	"ExprCalc/pkg/broker"
	"ExprCalc/pkg/config"
	"ExprCalc/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	config, err := config.NewConfig(".")
	if err != nil {
		panic(err)
	}

	logger := logger.NewLogger(config.Logger)

	server := server.NewServer(config.Server, logger)

	err = server.ConfigurateRedis(config.Redis, config.App)
	if err != nil {
		logger.Error("main: failed to open redis", zap.Error(err))
		return
	}
	exprService, err := expression.NewExpressionService(logger, config.Expr, broker.NewRabbit(logger, config.Rabbit))
	if err != nil {
		logger.Fatal("failed to start exprService", zap.Error(err))
	}
	defer exprService.Stop()

	rabbit := broker.NewRabbit(logger, config.Rabbit)
	err = rabbit.Ch.ExchangeDeclare(config.Expr.Exchange, "direct", true, false, false, false, nil)
	if err != nil {
		logger.Fatal("main.failed to declare exchange", zap.Error(err))
	}

	exprController := controllers.NewExpressionController(logger, config.Expr, rabbit, server.Redis)

	server.RegisterRouters([]controllers.Controller{exprController})

	// q, err := rabbit.Ch.QueueDeclare(config.Expr.ResultQueue, false, false, false, false, nil)
	// if err != nil {
	// 	logger.Fatal("main.failed to declare queue", zap.Error(err))
	// }

	// err = rabbit.Ch.QueueBind(q.Name, config.Expr.RouteKey, config.Expr.Exchange, false, nil)
	// if err != nil {
	// 	logger.Fatal("main.failed to bind queue", zap.Error(err))
	// }

	// out, err := rabbit.Ch.Consume(q.Name, "", false, false, false, false, nil)
	// if err != nil {
	// 	logger.Fatal("main.failed to consume queue", zap.Error(err))
	// }

	// body, _ := models.NewExpression("1+1").MarshalBinary()
	// corrID := gouid.Bytes(16)
	// err = rabbit.Ch.PublishWithContext(context.Background(), config.Expr.Exchange, config.Expr.RouteKey, false, false, amqp.Publishing{
	// 	ReplyTo:       q.Name,
	// 	CorrelationId: corrID.String(),
	// 	ContentType:   "application/json",
	// 	Body:          body,
	// })

	// fmt.Println("wait")
	// for msg := range out {
	// 	fmt.Println(string(msg.Body), msg.CorrelationId, msg.ReplyTo)
	// 	msg.Ack(false)
	// }

	server.Run()
}
