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
	exprController := controllers.NewExpressionController(logger, config.Expr, broker.NewRabbit(logger, config.Rabbit), server.Redis)

	exprService, err := expression.NewExpressionService(logger, config.Expr, broker.NewRabbit(logger, config.Rabbit), server.Redis)
	if err != nil {
		logger.Fatal("failed to start exprService", zap.Error(err))
	}
	defer exprService.Stop()

	server.RegisterRouters([]controllers.Controller{exprController})

	server.Run()
}
