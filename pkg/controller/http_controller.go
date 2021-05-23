package controller

import (
	"github.com/ereb-or-od/kenobi/pkg/controller/interfaces"
	"github.com/labstack/echo/v4"
)

type HttpController interface {
	interfaces.ControllerBase
	Endpoints() *map[string]map[string]echo.HandlerFunc
}
