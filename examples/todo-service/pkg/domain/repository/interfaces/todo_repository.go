package interfaces

import "github.com/ereb-or-od/kenobi/examples/todo-service/pkg/domain"

type TodoRepository interface {
	Create(todo *domain.Todo)
	FindById(id string) *domain.Todo
	Delete(id string)
}
