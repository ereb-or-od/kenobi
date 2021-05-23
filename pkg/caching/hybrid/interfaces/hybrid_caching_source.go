package interfaces

import (
	"context"
	"time"
)

type HybridCachingSource interface{
	GetOrSetValueByKey(ctx context.Context, key string, expiration time.Duration, callbackFunc func() (interface{}, error)) (interface{}, error)
	SetValue(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}
