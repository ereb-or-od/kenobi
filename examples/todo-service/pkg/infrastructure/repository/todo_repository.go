package repository

import (
	"github.com/ereb-or-od/kenobi/examples/todo-service/pkg/domain"
	"github.com/ereb-or-od/kenobi/examples/todo-service/pkg/domain/repository/interfaces"
	"github.com/ereb-or-od/kenobi/pkg/caching/inmemory"
	mem "github.com/ereb-or-od/kenobi/pkg/caching/inmemory/interfaces"
)

type todoRepository struct {
	database mem.InMemoryCachingSource
}

func (t todoRepository) Delete(id string) {
	t.database.DeleteValueByKey(id)
}

func (t todoRepository) FindById(id string) *domain.Todo {
	data := t.database.GetValueByKey(id)
	if data == nil {
		return new(domain.Todo)
	}
	todo := data.(*domain.Todo)
	return todo
}

func (t todoRepository) Create(todo *domain.Todo) {
	t.database.SetValue(todo.Id, todo)
}

func NewTodoRepository() interfaces.TodoRepository {
	inmemorySource, _ := inmemory.New()
	return &todoRepository{
		database: inmemorySource,
	}
}
