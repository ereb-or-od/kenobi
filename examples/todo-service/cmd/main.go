package main

import (
	"github.com/ereb-or-od/kenobi/examples/todo-service/pkg/api"
	"github.com/ereb-or-od/kenobi/pkg/server"
)

func main() {
	kenobiServer := server.New("todo_app").
		WithDefaultLogger().
		UseHttp().
		WithLoggingMiddleware().
		WithRecoverMiddleware().
		WithRequestIDMiddleware().
		WithAllowAnyCORSMiddleware().
		WithGzipMiddleware().
		WithHealthCheckMiddleware("/ping", "pong!").
		WithController(api.NewTodoController())
	kenobiServer.Start()

}
