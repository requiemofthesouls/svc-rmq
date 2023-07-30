package examplepb

import (
	"context"

	"github.com/requiemofthesouls/svc-rmq/client"
)

type EventsPublisherClient interface {
	UserRegistered(ctx context.Context, ev *UserRegisteredEvent) error
	UserLogged(ctx context.Context, ev *UserLoggedEvent) error
}

type eventsPublisherClient struct {
	client client.Manager
}

func NewEventsPublisherClient(client client.Manager) EventsPublisherClient {
	return &eventsPublisherClient{
		client: client,
	}
}

func (c *eventsPublisherClient) UserRegistered(ctx context.Context, msg *UserRegisteredEvent) error {
	return c.client.PublishToExchange(ctx, "events.user", "registered", msg)
}

func (c *eventsPublisherClient) UserLogged(ctx context.Context, msg *UserLoggedEvent) error {
	return c.client.PublishToExchange(ctx, "events.user", "logged", msg)
}
