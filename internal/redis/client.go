package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"order-service/internal/conf"
	"order-service/internal/constants"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

var ProviderSet = wire.NewSet(NewCache)

type RedisCache struct {
	client redis.UniversalClient
	ctx    context.Context
	log    *log.Helper
}

type RedisHandlerInterface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value []byte) error
}

func NewCache(conf *conf.Data, logger log.Logger) RedisHandlerInterface {
	l := log.NewHelper(logger)
	ctx := context.Background()

	conn := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        []string{conf.Redis.Addr},
		Password:     "",
		PoolSize:     int(conf.Redis.PoolSize),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	_, err := conn.Ping(ctx).Result()
	if err != nil {
		l.Errorf("Failed to connect to redis: %v", err)
		os.Exit(1)
	}
	l.Info("Connected to redis successfully!")

	return &RedisCache{
		client: conn,
		log:    l,
		ctx:    ctx,
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()

	if err != nil {
		errMsg := fmt.Sprintf(constants.RedisGetError, err)
		r.log.Error(errMsg)
		return "", errors.New(500, constants.RedisError, errMsg)
	}
	return val, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value []byte) error {
	err := r.client.Set(ctx, key, value, 0).Err()

	if err != nil {
		errMsg := fmt.Sprintf(constants.RedisSetError, err)
		r.log.Error(errMsg)
		return errors.New(500, constants.RedisError, errMsg)
	}
	return nil
}
