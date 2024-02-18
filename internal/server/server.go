package server

import (
	"ExprCalc/internal/server/controllers"
	"ExprCalc/pkg/config"
	"ExprCalc/pkg/repository/mongo"
	"ExprCalc/pkg/repository/redisdb"
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Server struct {
	Config      *config.ServerConfig
	Logger      *zap.Logger
	Router      *echo.Echo
	Controllers []controllers.Controller
	Redis       *redisdb.RedisDB
	Mongo       *mongo.MongoDB
}

func NewServer(config *config.ServerConfig, logger *zap.Logger) *Server {
	return &Server{
		Config:      config,
		Logger:      logger,
		Router:      echo.New(),
		Controllers: []controllers.Controller{},
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)
	s.Logger.Info("Starting server", zap.String("addr", addr))
	return s.Router.Start(addr)
}

func (s *Server) ConfigurateRedis(redisconfig *config.RedisDBConfig, appconfig *config.AppConfig) error {
	redis := redisdb.NewRedis(appconfig, redisconfig, s.Logger)

	err := redis.Open()
	if err != nil {
		s.Logger.Error("server.ConfigurateRedis: failed to open redis", zap.Error(err))
		return err
	}
	s.Redis = redis
	return nil
}

func (s *Server) RegisterRouters(routes []controllers.Controller) {
	for _, route := range routes {
		group := s.Router.Group(route.GetGroup())
		for _, middleware := range route.GetMiddleware() {
			group.Use(middleware)
		}
		for _, handler := range route.GetHandlers() {
			group.Add(handler.GetMethod(), handler.GetPath(), handler.GetHandler())
		}
	}
	s.Controllers = routes
}

/*
*	Пока не нужна.
* 	Есть какая-то проблема что при запуске выдает панику.
*	Так что будет зомби кодом.
 */
func (s *Server) ConfigurateMongo(config *config.MongoDBConfig) {
	m := mongo.New(config, s.Logger)

	err := m.Open()
	if err != nil {
		s.Logger.Fatal("server.ConfigurateMongo: error open mongo", zap.Error(err))
	}
	s.Mongo = m
}
