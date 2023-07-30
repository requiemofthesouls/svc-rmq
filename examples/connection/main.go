package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/requiemofthesouls/logger"

	"github.com/requiemofthesouls/svc-rmq/connection"
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
		Level:    "info",
		Encoding: "console",
		Caller:   true,
	}, nil); err != nil {
		fmt.Printf("logger.New error: %v\n", err)
		return
	}

	var logConn1 = log.With(logger.String("rmq.connection.name", "connection-1"))
	var rmqConn connection.Manager
	if rmqConn, err = connection.NewManager(logConn1, connection.Config{
		Host:     "localhost",
		Port:     5672,
		Username: "guest",
		Password: "guest",
		Params: connection.ConfigParams{
			ConnectionName: "svc-rmq-example-connection-connection-1",
		},
	}); err != nil {
		return
	}

	<-ctx.Done()
	rmqConn.Close()

	time.Sleep(time.Second * 3)
}
