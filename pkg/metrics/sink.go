package metrics

import (
	"fmt"
	"net/url"
)

// The MetricSink interface is used to transmit metrics information
// to an external system
type MetricSink interface {
	SetGauge(key []string, val float32)
	SetGaugeWithLabels(key []string, val float32, labels []Label)
	EmitKey(key []string, val float32)
	IncrCounter(key []string, val float32)
	IncrCounterWithLabels(key []string, val float32, labels []Label)
	AddSample(key []string, val float32)
	AddSampleWithLabels(key []string, val float32, labels []Label)
}

// EmptySink is used to just blackhole messages
type EmptySink struct{}

func (*EmptySink) SetGauge(key []string, val float32)                              {}
func (*EmptySink) SetGaugeWithLabels(key []string, val float32, labels []Label)    {}
func (*EmptySink) EmitKey(key []string, val float32)                               {}
func (*EmptySink) IncrCounter(key []string, val float32)                           {}
func (*EmptySink) IncrCounterWithLabels(key []string, val float32, labels []Label) {}
func (*EmptySink) AddSample(key []string, val float32)                             {}
func (*EmptySink) AddSampleWithLabels(key []string, val float32, labels []Label)   {}

// FanoutSink is used to sink to fanout values to multiple sinks
type FanoutSink []MetricSink

func (fh FanoutSink) SetGauge(key []string, val float32) {
	fh.SetGaugeWithLabels(key, val, nil)
}

func (fh FanoutSink) SetGaugeWithLabels(key []string, val float32, labels []Label) {
	for _, s := range fh {
		s.SetGaugeWithLabels(key, val, labels)
	}
}

func (fh FanoutSink) EmitKey(key []string, val float32) {
	for _, s := range fh {
		s.EmitKey(key, val)
	}
}

func (fh FanoutSink) IncrCounter(key []string, val float32) {
	fh.IncrCounterWithLabels(key, val, nil)
}

func (fh FanoutSink) IncrCounterWithLabels(key []string, val float32, labels []Label) {
	for _, s := range fh {
		s.IncrCounterWithLabels(key, val, labels)
	}
}

func (fh FanoutSink) AddSample(key []string, val float32) {
	fh.AddSampleWithLabels(key, val, nil)
}

func (fh FanoutSink) AddSampleWithLabels(key []string, val float32, labels []Label) {
	for _, s := range fh {
		s.AddSampleWithLabels(key, val, labels)
	}
}

// sinkURLFactoryFunc is an generic interface around the *SinkFromURL() function provided
// by each sink type
type sinkURLFactoryFunc func(*url.URL) (MetricSink, error)

// sinkRegistry supports the generic NewMetricSink function by mapping URL
// schemes to metric sink factory functions
var sinkRegistry = map[string]sinkURLFactoryFunc{
	"statsd":   NewStatsdSinkFromURL,
	"statsite": NewStatsiteSinkFromURL,
	"inmem":    NewInmemSinkFromURL,
}

// NewMetricSinkFromURL allows a generic URL input to configure any of the
// supported sinks. The scheme of the URL identifies the type of the sink, the
// and query parameters are used to set options.
//
// "statsd://" - Initializes a StatsdSink. The host and port are passed through
// as the "addr" of the sink
//
// "statsite://" - Initializes a StatsiteSink. The host and port become the
// "addr" of the sink
//
// "inmem://" - Initializes an InmemSink. The host and port are ignored. The
// "interval" and "duration" query parameters must be specified with valid
// durations, see NewInmemSink for details.
func NewMetricSinkFromURL(urlStr string) (MetricSink, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	sinkURLFactoryFunc := sinkRegistry[u.Scheme]
	if sinkURLFactoryFunc == nil {
		return nil, fmt.Errorf(
			"cannot create metric sink, unrecognized sink name: %q", u.Scheme)
	}

	return sinkURLFactoryFunc(u)
}
