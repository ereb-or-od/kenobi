package application

import (
	"context"
	"github.com/ereb-or-od/kenobi/pkg/mediator"
)

type FindTodoByIdQuery struct {
	Id string
}

type TodoContract struct {
	Id   string
	Name string
}

func (*FindTodoByIdQuery) Key() string { return "FindTodoByIdQuery" }

type FindTodoByIdQueryHandler struct {
	baseHandler *BaseHandler
}

func NewFindTodoByIdQueryHandler(baseHandler *BaseHandler) FindTodoByIdQueryHandler {
	return FindTodoByIdQueryHandler{baseHandler: baseHandler}
}

func (c FindTodoByIdQueryHandler) Handle(_ context.Context, query mediator.Message) (interface{}, error) {
	q := query.(*FindTodoByIdQuery)
	todo := c.baseHandler.repository.FindById(q.Id)
	return &TodoContract{
		Id:   todo.Id,
		Name: todo.Name,
	}, nil
}
