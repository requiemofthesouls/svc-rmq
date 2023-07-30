package def

import (
	"github.com/requiemofthesouls/container"
	logDef "github.com/requiemofthesouls/logger/def"

	interceptor2 "github.com/requiemofthesouls/svc-rmq/server/interceptor"
)

const (
	DIInterceptorRecovery = "rmq.server.middleware.recovery"
	DIInterceptorContext  = "rmq.server.middleware.context"
	DIInterceptorLogging  = "rmq.server.middleware.logging"
)

var List = []string{
	DIInterceptorContext,
	DIInterceptorLogging,
	DIInterceptorRecovery,
}

func init() {
	container.Register(func(builder *container.Builder, _ map[string]interface{}) error {
		return builder.Add(
			container.Def{
				Name: DIInterceptorRecovery,
				Build: func(cont container.Container) (interface{}, error) {
					var log logDef.Wrapper
					if err := cont.Fill(logDef.DIWrapper, &log); err != nil {
						return nil, err
					}

					return interceptor2.Recovery(log), nil
				},
			},
			container.Def{
				Name: DIInterceptorContext,
				Build: func(cont container.Container) (interface{}, error) {
					return interceptor2.Context(), nil
				},
			},
			container.Def{
				Name: DIInterceptorLogging,
				Build: func(cont container.Container) (interface{}, error) {
					var log logDef.Wrapper
					if err := cont.Fill(logDef.DIWrapper, &log); err != nil {
						return nil, err
					}

					return interceptor2.Logging(log), nil
				},
			},
		)
	})
}
