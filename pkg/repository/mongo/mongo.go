package mongo

import (
	"ExprCalc/pkg/config"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDB struct {
	Client *mongo.Client
	config *config.MongoDBConfig
	db     *mongo.Database
	logger *zap.Logger
}

func New(config *config.MongoDBConfig, logger *zap.Logger) *MongoDB {
	m := &MongoDB{
		config: config,
		logger: logger,
	}
	if m == nil {
		m.logger.Fatal("mongo.New: failed to create mongo")
	}
	return m
}

func (m *MongoDB) Open() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(m.config.URI))
	if err != nil {
		m.logger.Error("mongo.open: failed to connect to mongo", zap.Error(err))
		return err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		m.logger.Error("mongo.open: failed to ping mongo", zap.Error(err))
		return err
	}

	m.logger.Info("mongo.open: connected to mongo", zap.String("URI", m.config.URI))
	m.Client = client
	m.db = client.Database("mes")

	return nil
}

func (m *MongoDB) Close() error {
	err := m.Client.Disconnect(context.TODO())
	if err != nil {
		m.logger.Error("mongo.Close: failed to disconnect from mongo", zap.Error(err))
		return err
	}
	return nil
}
