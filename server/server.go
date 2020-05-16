package server

import (
	"context"
	"net"
	"net/http"
	"runtime"

	"github.com/rs/cors"
)

// https://github.com/soheilhy/cmux
// https://github.com/google/huproxy

func Run(ctx context.Context) {
	server := &http.Server{
		Addr: ":7532",
		BaseContext: func(lis net.Listener) context.Context {
			return ctx
		},
		ConnContext: func(ctx context.Context, conn net.Conn) context.Context {
			return ctx
		},
		Handler: cors.AllowAll().Handler(Handler()),
	}

	if runtime.GOOS == "linux" {
		server.Addr = ":80"
	}

	go func() {
		<-ctx.Done()
		server.Close()
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func Handler() http.HandlerFunc {

	return ServeContainerHTTP()
}
