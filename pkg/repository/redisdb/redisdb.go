package redisdb

import (
	"ExprCalc/pkg/config"
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisDB struct {
	Logger      *zap.Logger
	redisConfig *config.RedisDBConfig
	mutex       sync.RWMutex
	appConfig   *config.AppConfig
	Client      *redis.Client
}

func NewRedis(app_cfg *config.AppConfig, re_cfg *config.RedisDBConfig, logger *zap.Logger) *RedisDB {
	return &RedisDB{
		Logger:      logger,
		redisConfig: re_cfg,
		appConfig:   app_cfg,
	}
}

func (r *RedisDB) Open() error {
	options, err := redis.ParseURL(r.redisConfig.URI)
	if err != nil {
		r.Logger.Error("redis.Open: failed to parse redis uri", zap.Error(err))
		return err
	}

	client := redis.NewClient(options)
	if err := client.Ping(context.TODO()).Err(); err != nil {
		r.Logger.Error("redis.Open: failed to ping redis", zap.Error(err))
		return err
	}
	r.Client = client
	r.Logger.Info("redis.Open: connected to redis", zap.String("URI", r.redisConfig.URI))
	return nil
}

func (r *RedisDB) Close() error {
	return r.Client.Close()
}

func (r *RedisDB) WriteCache(ctx context.Context, key string, value interface{}) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	err := r.Client.Set(ctx, key, value, time.Duration(r.appConfig.CacheTTL)*time.Minute)
	return err.Err()
}

func (r *RedisDB) GetCache(ctx context.Context, key string) (interface{}, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.Client.Get(ctx, key).Result()
}

func (r *RedisDB) IsExist(ctx context.Context, key string) (int64, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.Client.Exists(ctx, key).Result()
}
