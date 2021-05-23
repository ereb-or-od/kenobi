package standalone

import (
	"context"
	logger "github.com/ereb-or-od/kenobi/pkg/logging/interfaces"
	marshallers "github.com/ereb-or-od/kenobi/pkg/marshalling/interfaces"
	"github.com/ereb-or-od/kenobi/pkg/redis/interfaces"
	"github.com/go-redis/redis/v8"
	"time"
)

type standaloneRedisServer struct {
	logger     logger.Logger
	client     *redis.Client
	marshaller marshallers.Marshaller
}

func (r standaloneRedisServer) DeleteValueByKey(ctx context.Context, key string) error {
	commandResult := r.client.Del(ctx, key)
	return commandResult.Err()
}

func (r standaloneRedisServer) GetValueByKey(ctx context.Context, key string, result interface{}) error {
	commandResult := r.client.Get(ctx, key)
	if commandResult.Err() != nil {
		if commandResult.Err().Error() == "redis: nil" {
			return nil
		}
		return commandResult.Err()
	}

	resultAsByteArray, err := commandResult.Bytes()
	if err != nil {
		return err
	}
	err = r.marshaller.Unmarshall(resultAsByteArray, &result)
	if err != nil {
		return err
	}

	return nil
}

func (r standaloneRedisServer) SetValue(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	byteArray, err := r.marshaller.Marshall(&value)
	if err != nil {
		return err
	}
	commandResult := r.client.Set(ctx, key, byteArray, expiration)
	if commandResult.Err() != nil {
		return commandResult.Err()
	}
	return nil
}

func New(logger logger.Logger, marshaller marshallers.Marshaller, options *StandaloneRedisServerOptions) interfaces.RedisServer {
	rdb := redis.NewClient(&redis.Options{
		Addr:     options.Address,
		Password: options.Password,
	})

	return &standaloneRedisServer{
		logger:     logger,
		client:     rdb,
		marshaller: marshaller,
	}
}
