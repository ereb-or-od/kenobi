package middleware

import (
	"context"
	"github.com/ereb-or-od/kenobi/pkg/rabbitmq/consumer"
	"github.com/streadway/amqp"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "amqpextra/middleware context value " + k.name
}

func nack(ctx context.Context, msg amqp.Delivery) interface{} {
	if err := msg.Nack(false, false); err != nil {
	}
	return nil
}

func wrap(fn func(ctx context.Context, msg amqp.Delivery, next consumer.Handler) interface{}) func(next consumer.Handler) consumer.Handler {
	return func(next consumer.Handler) consumer.Handler {
		return consumer.HandlerFunc(func(ctx context.Context, msg amqp.Delivery) interface{} {
			return fn(ctx, msg, next)
		})
	}
}
