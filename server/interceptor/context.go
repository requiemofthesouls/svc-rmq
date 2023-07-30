package interceptor

import (
	"context"

	userclient "github.com/requiemofthesouls/user-client"

	"github.com/requiemofthesouls/svc-rmq/internal"
	"github.com/requiemofthesouls/svc-rmq/server/consumer"
)

func Context() consumer.Interceptor {
	return func(ctx context.Context, msg *consumer.Message, handler consumer.Handler) error {
		ctx = userclient.RequestIDToContext(ctx, msg.MessageId)

		if requestID, ok := msg.Headers[internal.HeaderRequestID].(string); ok {
			ctx = userclient.RequestIDToContext(ctx, requestID)
		}

		if userClient, ok := msg.Headers[internal.HeaderUserClient].(string); ok {
			ctx = userclient.ClientToContext(ctx, userclient.UnmarshallClient([]byte(userClient)))
		}

		return handler(ctx, msg)
	}
}
