package failover

import "github.com/ereb-or-od/kenobi/pkg/redis/options"

type RedisServerOptions struct {
	MasterName       string
	SentinelAddrs    []string
	SentinelPassword string
	*options.RedisServerOptions
}
