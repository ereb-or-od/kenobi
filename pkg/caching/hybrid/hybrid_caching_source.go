package hybrid

import (
	"context"
	ds "github.com/ereb-or-od/kenobi/pkg/caching/distributed/interfaces"
	"github.com/ereb-or-od/kenobi/pkg/caching/hybrid/interfaces"
	in "github.com/ereb-or-od/kenobi/pkg/caching/inmemory/interfaces"
	"time"
)

type hybridCachingSource struct {
	inmemoryCachingSource    in.InMemoryCachingSource
	distributedCachingSource ds.DistributedCachingSource
}

func (h hybridCachingSource) GetOrSetValueByKey(ctx context.Context, key string, expiration time.Duration, callbackFunc func() (interface{}, error)) (interface{}, error) {
	if dataInInMemoryCacheStore, err := h.inmemoryCachingSource.GetValueByKey(key); err != nil {
		return nil, err
	} else {
		if dataInInMemoryCacheStore == nil {
			var dataInDistributedCacheStore interface{}
			if err = h.distributedCachingSource.GetValueByKey(ctx, key, &dataInDistributedCacheStore); err != nil {
				return nil, err
			} else {
				if dataInDistributedCacheStore == nil {
					if callbackResult, err := callbackFunc(); err != nil {
						return nil, err
					} else {
						if err = h.SetValue(ctx, key, callbackResult, expiration); err != nil {
							return nil, err
						}
						return callbackResult, nil
					}
				} else {
					if err = h.inmemoryCachingSource.SetValue(key, dataInDistributedCacheStore); err != nil {
						return nil, err
					}
					return dataInDistributedCacheStore, nil
				}
			}
		} else {
			return dataInInMemoryCacheStore, nil
		}
	}
}

func (h hybridCachingSource) SetValue(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := h.inmemoryCachingSource.SetValue(key, value); err != nil {
		return err
	} else {
		if err = h.distributedCachingSource.SetValue(ctx, key, value, expiration); err != nil {
			return err
		}
		return nil
	}
}

func New(inmemoryCacheSource in.InMemoryCachingSource, distributedCacheSource ds.DistributedCachingSource) interfaces.HybridCachingSource {
	return &hybridCachingSource{
		inmemoryCachingSource:    inmemoryCacheSource,
		distributedCachingSource: distributedCacheSource,
	}
}
