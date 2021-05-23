package api

import (
	"context"
	"github.com/ereb-or-od/kenobi/examples/todo-service/pkg/application"
	"github.com/ereb-or-od/kenobi/pkg/controller"
	"github.com/ereb-or-od/kenobi/pkg/mediator"
	"github.com/labstack/echo/v4"
)

type TodoController struct {
	mediator *mediator.Mediator
}

func (h TodoController) Name() string {
	return "todo"
}

func (h TodoController) Prefix() string {
	return "todo"
}

func (h TodoController) Version() string {
	return "v1"
}

func (h TodoController) Endpoints() *map[string]map[string]echo.HandlerFunc {
	return &map[string]map[string]echo.HandlerFunc{
		"": {
			"POST": func(echoContext echo.Context) error {
				command := new(application.CreateTodoCommand)
				if err := echoContext.Bind(command); err != nil {
					return err
				}
				result, _ := h.mediator.Send(context.Background(), command)
				return echoContext.JSON(200, result)
			},
		},
		"/:id": {
			"GET": func(echoContext echo.Context) error {
				result, _ := h.mediator.Send(context.Background(), &application.FindTodoByIdQuery{Id: echoContext.Param("id")})
				return echoContext.JSON(200, result)
			},
			"DELETE": func(echoContext echo.Context) error {
				result, _ := h.mediator.Send(context.Background(), &application.DeleteTodoByIdCommand{Id: echoContext.Param("id")})
				return echoContext.JSON(200, result)
			},
		},
	}
}

func NewTodoController() controller.HttpController {
	baseHandler := application.NewBaseHandler()
	m, _ := mediator.NewContext().
		RegisterHandler(&application.CreateTodoCommand{}, application.NewCreateTodoCommandHandler(baseHandler)).
		RegisterHandler(&application.FindTodoByIdQuery{}, application.NewFindTodoByIdQueryHandler(baseHandler)).
		RegisterHandler(&application.DeleteTodoByIdCommand{}, application.NewDeleteTodoByIdCommandHandler(baseHandler)).
		Build()

	return &TodoController{
		mediator: m,
	}
}
