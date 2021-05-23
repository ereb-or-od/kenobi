package middleware

import (
	"context"
	"github.com/ereb-or-od/kenobi/pkg/rabbitmq/consumer"

	"github.com/streadway/amqp"
)

func HasCorrelationID() consumer.Middleware {
	return wrap(func(ctx context.Context, msg amqp.Delivery, next consumer.Handler) interface{} {
		if msg.CorrelationId == "" {
			return nack(ctx, msg)
		}

		return next.Handle(ctx, msg)
	})
}
