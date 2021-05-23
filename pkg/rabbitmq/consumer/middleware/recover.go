package middleware

import (
	"context"
	"github.com/ereb-or-od/kenobi/pkg/rabbitmq/consumer"
	"github.com/streadway/amqp"
)

func Recover() consumer.Middleware {
	return wrap(func(ctx context.Context, msg amqp.Delivery, next consumer.Handler) (result interface{}) {
		defer func() {
			if e := recover(); e != nil {
				if nackErr := msg.Nack(false, false); nackErr != nil {
				}
			}
		}()

		return next.Handle(ctx, msg)
	})
}
