package main

import (
	"ExprCalc/internal/server"
	"ExprCalc/internal/server/controllers"
	"ExprCalc/internal/services/expression"
	"ExprCalc/pkg/broker"
	"ExprCalc/pkg/config"
	"ExprCalc/pkg/logger"
	"context"
	"fmt"

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

	exprService, err := expression.NewExpressionService(logger, config.Expr, broker.NewRabbit(logger, config.Rabbit), server.Redis)
	if err != nil {
		logger.Fatal("failed to start exprService", zap.Error(err))
	}
	defer exprService.Stop()

	res, err := server.Redis.GetAllKeysByPattern(context.TODO(), "worker")
	fmt.Println(res)
	exprController := controllers.NewExpressionController(logger, config.Expr, broker.NewRabbit(logger, config.Rabbit), server.Redis)

	server.RegisterRouters([]controllers.Controller{exprController})

	server.Run()
}
