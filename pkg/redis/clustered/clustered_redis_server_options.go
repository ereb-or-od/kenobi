package clustered

import "github.com/ereb-or-od/kenobi/pkg/redis/options"

type RedisServerOptions struct {
	Addresses []string
	*options.RedisServerOptions
}
