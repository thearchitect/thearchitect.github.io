package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"runtime"

	"github.com/rs/cors"
	"github.com/thearchitect/thearchitect.github.io/server/resources"
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
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", ServeContainerHTTP())

	index, webapp := resources.IndexHTML(true)

	webapp.Mount(mux)

	mux.Handle("/", index)

	return func(w http.ResponseWriter, q *http.Request) {

		log.Println("request", q.URL.String())

		mux.ServeHTTP(w, q)
	}
}
