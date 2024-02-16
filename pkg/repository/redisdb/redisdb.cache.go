package redisdb

import (
	"context"
	"time"
)

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

func (r *RedisDB) GetAllKeysByPattern(ctx context.Context, pattern string) ([]string, error) {
	return r.Client.Keys(ctx, pattern).Result()
}
