package interfaces

import (
	"context"
	"time"
)

type RedisServer interface{
	GetValueByKey(ctx context.Context, key string, result interface{})  error
	SetValue(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	DeleteValueByKey(ctx context.Context, key string) error
}
