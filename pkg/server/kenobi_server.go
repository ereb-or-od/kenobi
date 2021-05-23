package server

import (
	"fmt"
	controllers "github.com/ereb-or-od/kenobi/pkg/controller"
	controllerBase "github.com/ereb-or-od/kenobi/pkg/controller/interfaces"
	"github.com/ereb-or-od/kenobi/pkg/http/middlewares"
	"github.com/ereb-or-od/kenobi/pkg/logging"
	logger "github.com/ereb-or-od/kenobi/pkg/logging/interfaces"
	"github.com/ereb-or-od/kenobi/pkg/logging/options"
	serverOption "github.com/ereb-or-od/kenobi/pkg/server/options"
	"github.com/ereb-or-od/kenobi/pkg/utilities"
	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	apmecho "github.com/opentracing-contrib/echo"
	"github.com/opentracing/opentracing-go"
	"github.com/tylerb/graceful"
	"time"
)

var (
	defaultServerPort                             = 80
	defaultExcludedEndpointsFromMetricsAndTracing = []string{"/metrics", "/stats", "/_stats", "/ping", "/health-check", "/healthy"}
)

type KenobiServer struct {
	serverOptions *serverOption.KenobiServerOptions
	http          *echo.Echo
	controller    controllerBase.ControllerBase
	logger        logger.Logger
}

func New(name string) *KenobiServer {
	return &KenobiServer{
		serverOptions: &serverOption.KenobiServerOptions{Name: name},
	}
}

func (k *KenobiServer) WithLogger(logger logger.Logger) *KenobiServer {
	if logger == nil {
		panic("logger must be specified")
	}
	if k.logger != nil {
		panic("you can not register logger more than once")
	}
	k.logger = logger
	return k
}

func (k *KenobiServer) WithDefaultLogger(options ...*options.LoggerOptions) *KenobiServer {
	if k.logger != nil {
		panic("you can not register logger more than once")
	}
	if defaultLogger, err := logging.New(options...); err != nil {
		panic(err)
	} else {
		k.logger = defaultLogger
	}

	return k
}

func (k *KenobiServer) UseHttp() *KenobiServer {
	k.http = echo.New()
	k.http.HideBanner = true
	return k
}

func (k *KenobiServer) UsePrometheus(excludedEndpoints ...string) *KenobiServer {
	k.serverOptions.Metric = &serverOption.KenobiServerMetricOptions{ExcludedEndpoints: excludedEndpoints}
	p := prometheus.NewPrometheus(k.serverOptions.Name, k.defaultEndpointSkipper)
	p.Use(k.http)

	return k
}

func (k *KenobiServer) UseOpenTracing() *KenobiServer {
	opentracing.SetGlobalTracer(opentracing.GlobalTracer())
	k.http.Use(apmecho.Middleware(k.serverOptions.Name))
	return k
}

func (k *KenobiServer) defaultEndpointSkipper(c echo.Context) bool {
	if utilities.ContainsInStringSlice(defaultExcludedEndpointsFromMetricsAndTracing, c.Path()) {
		return true
	} else {
		return false
	}
}

func (k *KenobiServer) WithLoggingMiddleware() *KenobiServer {
	k.http.Use(middlewares.LoggingMiddleware(k.logger))

	return k
}
func (k *KenobiServer) WithGzipMiddleware() *KenobiServer {
	k.http.Use(middleware.Gzip())
	return k
}

func (k *KenobiServer) WithRequestIDMiddleware() *KenobiServer {
	k.http.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{Generator: func() string {
		return uuid.NewString()
	}}))
	return k
}

func (k *KenobiServer) WithAllowAnyCORSMiddleware() *KenobiServer {
	k.http.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"*"},
	}))
	return k
}

func (k *KenobiServer) WithCORSMiddleware(allowsOrigins []string, allowsHeaders []string, allowsMethods []string) *KenobiServer {
	k.http.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowsOrigins,
		AllowHeaders: allowsHeaders,
		AllowMethods: allowsMethods,
	}))
	return k
}

func (k *KenobiServer) WithTimeoutMiddleware(duration time.Duration) *KenobiServer {
	k.http.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: duration,
	}))
	return k
}

func (k *KenobiServer) WithHealthCheckMiddleware(path string, response string) *KenobiServer {
	k.http.Use(middlewares.HealthCheckMiddleware(path, response))
	return k
}

func (k *KenobiServer) WithRecoverMiddleware() *KenobiServer {
	k.http.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         4 << 10,
		DisableStackAll:   false,
		DisablePrintStack: false,
		LogLevel:          log.ERROR,
	}))
	return k
}

func (k *KenobiServer) WithCustomMiddlewares(middlewares ...echo.MiddlewareFunc) *KenobiServer {
	if middlewares == nil || len(middlewares) == 0 {
		panic("middlewares must be specified")
	}
	for _, customMiddleware := range middlewares {
		k.http.Use(customMiddleware)
	}
	return k
}

func (k *KenobiServer) WithNewRelicMiddleware(licenceKey string) *KenobiServer {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(k.serverOptions.Name),
		newrelic.ConfigLicense(licenceKey))
	if err != nil {
		panic(err)
	}
	k.http.Use(nrecho.Middleware(app))
	return k
}
func (k *KenobiServer) WithController(controller controllerBase.ControllerBase) *KenobiServer {
	if controller == nil {
		panic("controller must be specified")
	}
	k.controller = controller

	var (
		httpController controllers.HttpController
	)
	switch svc := controller.(type) {
	case controllers.HttpController:
		httpController = svc
	default:
		panic("controller must implement to Service interface")
	}

	if httpController != nil {
		for path, endpintMethod := range *httpController.Endpoints() {
			for method, endpointHandler := range endpintMethod {
				var endpoint string
				if len(httpController.Prefix()) > 0 {
					endpoint += fmt.Sprintf("/%s", httpController.Prefix())
				}
				if len(httpController.Version()) > 0 {
					endpoint += fmt.Sprintf("/%s", httpController.Version())
				}
				endpoint += fmt.Sprintf("%s", path)
				httpController.Version()
				k.http.Add(method, endpoint, endpointHandler)
			}
		}
	}
	return k
}

func (k *KenobiServer) Start() {
	k.http.Server.Addr = fmt.Sprintf(":%d", defaultServerPort)
	k.http.Logger.Fatal(graceful.ListenAndServe(k.http.Server, 5*time.Second))
}

func (k *KenobiServer) StartWithOptions(options *serverOption.KenobiServerStartOptions) {
	port := fmt.Sprintf(":%d", options.Port)
	if options.GracefullyShutdown {
		k.http.Server.Addr = port
		k.http.Logger.Fatal(graceful.ListenAndServe(k.http.Server, options.GracefullyShutdownTimeoutPeriod))
	} else {
		k.http.Logger.Fatal(k.http.Start(port))
	}
}
