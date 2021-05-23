package middleware

import (
	"context"
	"github.com/ereb-or-od/kenobi/pkg/rabbitmq/consumer"
	"github.com/streadway/amqp"
)

func HasReplyTo() consumer.Middleware {
	return wrap(func(ctx context.Context, msg amqp.Delivery, next consumer.Handler) interface{} {
		if msg.ReplyTo == "" {
			return nack(ctx, msg)
		}

		return next.Handle(ctx, msg)
	})
}
