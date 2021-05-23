package domain

import (
	"github.com/google/uuid"
	"time"
)

type Todo struct {
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(name string) *Todo {
	return &Todo{
		Id:        uuid.NewString(),
		Name:      name,
		CreatedAt: time.Now(),
	}
}
