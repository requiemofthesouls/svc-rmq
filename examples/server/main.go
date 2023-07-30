package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/requiemofthesouls/logger"

	"github.com/requiemofthesouls/svc-rmq/connection"
	"github.com/requiemofthesouls/svc-rmq/examples/server/listener"
	examplepb "github.com/requiemofthesouls/svc-rmq/examples/server/pb"
	"github.com/requiemofthesouls/svc-rmq/server"
	"github.com/requiemofthesouls/svc-rmq/server/consumer"
	interceptor2 "github.com/requiemofthesouls/svc-rmq/server/interceptor"
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
		Level:    "debug",
		Encoding: "json",
		Caller:   true,
	}, nil); err != nil {
		fmt.Printf("logger.New error: %v\n", err)
		return
	}

	var rmqConn connection.Manager
	if rmqConn, err = connection.NewManager(log, connection.Config{
		Host:     "localhost",
		Port:     5672,
		Username: "guest",
		Password: "guest",
		Params: connection.ConfigParams{
			ConnectionName: "svc-rmq-example-connection-connection-1",
		},
	}); err != nil {
		log.Error("connection.NewManager error", logger.Error(err))
		return
	}

	log.Info("connection started")

	var kServer server.Manager
	if kServer, err = server.NewManager(
		log,
		rmqConn,
		[]server.ListenerRegistrant{
			func(srv server.Manager) {
				examplepb.RegisterCommonEventsServer(srv, listener.NewCommonEventsListener())
			},
		},
		[]consumer.Interceptor{interceptor2.Recovery(log), interceptor2.Context(), interceptor2.Logging(log)},
		map[string]*server.QueueHandlerSettings{
			"h.svc-rmq.UserLoggedHandler": {
				QOS:                    15,
				NumConsumers:           15,
				DelayNumFailedAttempts: 15,
				DelayTTL:               2000,
			},
		},
	); err != nil {
		log.Error("server.NewManager error", logger.Error(err))
		return
	}

	go func(ctx context.Context) {
		<-ctx.Done()
		log.Info("context done")

		kServer.CloseAll()
		log.Info("server closed")
	}(ctx)

	kServer.StartAll()
}
