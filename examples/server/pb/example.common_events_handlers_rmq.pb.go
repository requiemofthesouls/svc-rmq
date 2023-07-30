package examplepb

import (
	"context"

	"github.com/requiemofthesouls/svc-rmq/server"
	"github.com/requiemofthesouls/svc-rmq/server/consumer"
)

type CommonEventsServer interface {
	UserRegisteredHandler(context.Context, *UserRegisteredEvent) error
	UserLoggedHandler(context.Context, *UserLoggedEvent) error
}

func RegisterCommonEventsServer(s server.Manager, srv CommonEventsServer) {
	s.RegisterService(
		consumer.MapQueueHandlers{
			"h.svc-rmq.UserRegisteredHandler": consumer.QueueHandlerItem{
				ExchangeName: "events.user",
				RoutingKey:   "registered",
				Handler:      NewCommonEventsHandlersUserRegisteredHandler(srv),
			},
			"h.svc-rmq.UserLoggedHandler": consumer.QueueHandlerItem{
				ExchangeName: "events.user",
				RoutingKey:   "logged",
				Handler:      NewCommonEventsHandlersUserLoggedHandler(srv),
			},
		},
	)
}

type CommonEventsHandlersUserRegisteredHandler struct {
	listener CommonEventsServer
}

func NewCommonEventsHandlersUserRegisteredHandler(listener CommonEventsServer) CommonEventsHandlersUserRegisteredHandler {
	return CommonEventsHandlersUserRegisteredHandler{
		listener: listener,
	}
}

func (h CommonEventsHandlersUserRegisteredHandler) Handle(ctx context.Context, dec func(message interface{}) error) error {
	var ev UserRegisteredEvent
	if err := dec(&ev); err != nil {
		return err
	}

	return h.listener.UserRegisteredHandler(ctx, &ev)
}

type CommonEventsHandlersUserLoggedHandler struct {
	listener CommonEventsServer
}

func NewCommonEventsHandlersUserLoggedHandler(listener CommonEventsServer) CommonEventsHandlersUserLoggedHandler {
	return CommonEventsHandlersUserLoggedHandler{
		listener: listener,
	}
}

func (h CommonEventsHandlersUserLoggedHandler) Handle(ctx context.Context, dec func(message interface{}) error) error {
	var ev UserLoggedEvent
	if err := dec(&ev); err != nil {
		return err
	}

	if err := h.listener.UserLoggedHandler(ctx, &ev); err != nil {
		return err
	}

	return nil
}
