package middlewares

import (
	"fmt"
	logger "github.com/ereb-or-od/kenobi/pkg/logging/interfaces"
	"github.com/ereb-or-od/kenobi/pkg/utilities"
	"github.com/labstack/echo/v4"
	"strings"
	"time"
)

func LoggingMiddleware(logger logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}

			userIp, _ := utilities.GetIP(req)
			logParameters := map[string]interface{}{
				"remote_ip":          c.RealIP(),
				"latency":            time.Since(start).String(),
				"request.host":       req.Host,
				"request.method":     req.Method,
				"request.uri":        req.RequestURI,
				"response.status":    res.Status,
				"response.size":      res.Size,
				"request.user_agent": req.UserAgent(),
				"user.ip":            userIp,
				"forwarded-ip":       utilities.GetForwardedIP(req),
			}
			var headers string
			if req.Header != nil {
				for key, value := range req.Header {
					headers = fmt.Sprintf("key:%s value:%s", key, strings.Join(value, ";"))
				}
				logParameters["http:headers"] = headers
			}

			var formParameters string
			if req.Form != nil {
				for key, value := range req.Form {
					headers = fmt.Sprintf("key:%s value:%s", key, strings.Join(value, ";"))
				}
				logParameters["http:form"] = formParameters
			}

			if req.Cookies() != nil {
				for _, cookie := range req.Cookies() {
					logParameters[cookie.Name] = cookie.String()
				}
			}
			n := res.Status
			switch {
			case n >= 500:
				logger.Error("[server-http-log]", err, logParameters)
			case n >= 400:
				logger.Warn("[server-http-log]", logParameters)
			case n >= 300:
				logger.Warn("[server-http-log]", logParameters)
			default:
				logger.Info("[server-http-log]", logParameters)
			}
			return
		}
	}
}
