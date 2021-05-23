package redis

import (
	"context"
	"github.com/ereb-or-od/kenobi/pkg/caching/distributed/interfaces"
	redis "github.com/ereb-or-od/kenobi/pkg/redis/interfaces"
	"time"
)

type redisServerCachingSource struct {
	redisServer redis.RedisServer
}

func (r redisServerCachingSource) DeleteValueByKey(ctx context.Context, key string) error {
	return r.redisServer.DeleteValueByKey(ctx, key)
}

func (r redisServerCachingSource) GetValueByKey(ctx context.Context, key string, result interface{}) error {
	return r.redisServer.GetValueByKey(ctx, key, result)
}

func (r redisServerCachingSource) SetValue(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.redisServer.SetValue(ctx, key, value, expiration)
}

func New(redisServer redis.RedisServer) interfaces.DistributedCachingSource {
	return &redisServerCachingSource{
		redisServer: redisServer,
	}
}
