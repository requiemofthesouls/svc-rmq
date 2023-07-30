package listener

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/requiemofthesouls/logger"

	examplepb "github.com/requiemofthesouls/svc-rmq/examples/server/pb"
)

func NewCommonEventsListener() examplepb.CommonEventsServer {
	return &listener{}
}

type (
	listener struct {
	}
)

func (l *listener) UserRegisteredHandler(ctx context.Context, _ *examplepb.UserRegisteredEvent) error {
	getLogger(ctx).Info("UserRegisteredHandler")

	return errors.New("UserRegisteredHandler error")
}

func (l *listener) UserLoggedHandler(ctx context.Context, _ *examplepb.UserLoggedEvent) error {
	getLogger(ctx).Info("UserLoggedHandler")
	return nil
}

func getLogger(ctx context.Context) logger.Wrapper {
	return logger.NewFromZap(ctxzap.Extract(ctx))
}
