package mediator


import "context"

type (
	Sender interface {
		Send(context.Context, Message)  (interface{}, error)
	}
	Builder interface {
		RegisterHandler(request Message, handler RequestHandler) Builder
		UseBehaviour(PipelineBehaviour) Builder
		Use(fn func(context.Context, Message, Next)  (interface{}, error)) Builder
		Build() (*Mediator, error)
	}
	RequestHandler interface {
		Handle(context.Context, Message) (interface{}, error)
	}
	PipelineBehaviour interface {
		Process(context.Context, Message, Next)  (interface{}, error)
	}
	Message interface {
		Key() string
	}
)
