package application

import (
	"context"
	"github.com/ereb-or-od/kenobi/pkg/mediator"
)

type DeleteTodoByIdCommand struct {
	Id string
}

func (*DeleteTodoByIdCommand) Key() string { return "DeleteTodoByIdCommand" }

type DeleteTodoByIdCommandHandler struct {
	baseHandler *BaseHandler
}

func NewDeleteTodoByIdCommandHandler(baseHandler *BaseHandler) DeleteTodoByIdCommandHandler {
	return DeleteTodoByIdCommandHandler{baseHandler: baseHandler}
}

func (c DeleteTodoByIdCommandHandler) Handle(_ context.Context, command mediator.Message) (interface{}, error) {
	cmd := command.(*DeleteTodoByIdCommand)
	c.baseHandler.repository.Delete(cmd.Id)
	return nil, nil
}
