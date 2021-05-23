package options

import "time"

type KenobiServerOptions struct {
	Name    string
	Metric  *KenobiServerMetricOptions
}

type KenobiServerMetricOptions struct{
	ExcludedEndpoints []string
}

type KenobiServerStartOptions struct {
	Port                            int
	GracefullyShutdown              bool
	GracefullyShutdownTimeoutPeriod time.Duration
}

type KenobiServerJeagerOptions struct {
	AgentHost   string
	AgentPort   string
	Endpoint    string
	User        string
	Password    string
}