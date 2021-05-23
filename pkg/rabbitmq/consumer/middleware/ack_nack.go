package middleware

import (
	"context"
	"github.com/ereb-or-od/kenobi/pkg/rabbitmq/consumer"

	"github.com/streadway/amqp"
)

const Ack = "ack"
const Nack = "nack_requeue"
const Requeue = "requeue"

func AckNack() consumer.Middleware {
	return func(next consumer.Handler) consumer.Handler {
		fn := func(ctx context.Context, msg amqp.Delivery) interface{} {
			result := next.Handle(ctx, msg)
			if result == nil {
				return nil
			}

			switch result {
			case Ack:
				if err := msg.Ack(false); err != nil {
					return nil
				}
				return nil
			case Nack:
				if err := msg.Nack(false, false); err != nil {
					return nil
				}
				return nil
			case Requeue:
				if err := msg.Nack(false, true); err != nil {
					return nil
				}
				return nil
			}

			return result
		}

		return consumer.HandlerFunc(fn)
	}
}
