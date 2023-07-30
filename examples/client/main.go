package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/requiemofthesouls/logger"
	userclient "github.com/requiemofthesouls/user-client"
	userclientpb "github.com/requiemofthesouls/user-client/pb"

	"github.com/requiemofthesouls/svc-rmq/client"
	"github.com/requiemofthesouls/svc-rmq/connection"
	examplepb "github.com/requiemofthesouls/svc-rmq/examples/client/pb"
)

func main() {
	ctx, ctxCancel := context.WithCancel(context.Background())

	// graceful stop
	go func(cancelFunc context.CancelFunc) {
		var c = make(chan os.Signal, 1)
		signal.Notify(c,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
		)

		<-c

		cancelFunc()
	}(ctxCancel)

	var (
		log logger.Wrapper
		err error
	)
	if log, err = logger.New(logger.Config{
		Level:      "debug",
		Encoding:   "console",
		Caller:     true,
		Stacktrace: "error",
	}, nil); err != nil {
		fmt.Printf("logger.New error: %v\n", err)
		return
	}

	var kConn connection.Manager
	if kConn, err = connection.NewManager(log, connection.Config{
		Host:     "localhost",
		Port:     5672,
		Username: "guest",
		Password: "guest",
		Params: connection.ConfigParams{
			ConnectionName: "example-rmq-client-connection-1",
			Heartbeat:      10,
			Locale:         "en_US",
		},
	}); err != nil {
		log.Error("connection.NewManager error", logger.Error(err))
		return
	}

	log.Info("connection started")

	var kClient client.Manager
	if kClient, err = client.NewManager(
		log,
		kConn,
		nil,
	); err != nil {
		log.Error("client.NewManager error", logger.Error(err))
		return
	}

	eventsPublisherClient := examplepb.NewEventsPublisherClient(kClient)

	ticker := time.NewTicker(time.Second * 2)

	for {
		select {
		case <-ticker.C:
			if err := publishEvents(eventsPublisherClient); err != nil {
				log.Error("publishEvents failed", logger.Error(err))
			}
		case <-ctx.Done():
			kClient.Close()

			fmt.Println("connection finished")
			return
		}
	}
}

func publishEvents(eventsPublisherClient examplepb.EventsPublisherClient) error {
	ctxCommon := userclient.ClientToContext(
		context.Background(),
		&userclientpb.Client{
			Ip:       "192.168.0.1",
			Host:     "localhost",
			Language: "en",
			Location: "RUS",
			Platform: "web",
		},
	)

	if err := eventsPublisherClient.UserRegistered(
		userclient.RequestIDToContext(ctxCommon, userclient.GenerateRequestID()),
		&examplepb.UserRegisteredEvent{
			UserId:       1,
			FirstName:    "",
			LastName:     "",
			Email:        "",
			Country:      "",
			RegisteredAt: time.Now().Unix(),
		},
	); err != nil {
		return fmt.Errorf("eventsPublisherClient.UserRegistered error: %v", err)
	}

	if err := eventsPublisherClient.UserLogged(
		userclient.RequestIDToContext(ctxCommon, userclient.GenerateRequestID()),
		&examplepb.UserLoggedEvent{
			UserId:    1,
			FirstName: "",
			LastName:  "",
			Email:     "",
			Country:   "",
			LoggedAt:  time.Now().Unix(),
		},
	); err != nil {
		return fmt.Errorf("eventsPublisherClient.UserLogged error: %v", err)
	}

	return nil
}
