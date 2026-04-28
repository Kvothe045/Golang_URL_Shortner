package store

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
	mu     sync.RWMutex //not needed with Redis, but we do it just for the sake of learning
}

func Newrdb() *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	return &RedisStore{
		client: rdb,
		ctx:    ctx,
	}
}

func (rs *RedisStore) Save(url string, code string) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.client.Set(rs.ctx, url, code, 60*time.Second)
	rs.client.Set(rs.ctx, code, url, 60*time.Second)
}

func (rs *RedisStore) GetEncodedURL(url string) (string, bool) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	code, err := rs.client.Get(rs.ctx, url).Result()

	return code, err == nil
}

func (rs *RedisStore) GetOriginalURL(code string) (string, bool) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	url, err := rs.client.Get(rs.ctx, code).Result()

	return url, err == nil
}

func (rs *RedisStore) IsHealthy() bool {
	err := rs.client.Ping(rs.ctx).Err()
	return err == nil
}
