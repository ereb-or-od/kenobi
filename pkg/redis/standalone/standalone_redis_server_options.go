package standalone

import "github.com/ereb-or-od/kenobi/pkg/redis/options"

type StandaloneRedisServerOptions struct {
	Address string
	*options.RedisServerOptions
}
