package main

import (
	"context"
	"github.com/thearchitect/thearchitect.github.io/server"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/thearchitect/thearchitect.github.io/internal/docker"
	"github.com/thearchitect/thearchitect.github.io/ui"
)

// https://github.com/containers/libpod
// https://docs.docker.com/engine/api/v1.24/#attach-to-a-container-websocket
// https://github.com/fogleman/pt
// https://github.com/cmatsuoka/figlet
// docker run --rm -it geertjohan/gomatrix

//go:generate go generate ./server/resources

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	{
		var sch = make(chan os.Signal)
		signal.Notify(sch, os.Interrupt, os.Kill, syscall.SIGHUP)

		go func() {
			log.Println("SIGNAL", <-sch)
			cancel()
		}()
	}

	if inDocker, err := docker.WhereAmI(); err != nil {
		panic(err)
	} else {

		if inDocker {
			docker.Hide()

			ui.Run(ctx)
		} else if runtime.GOOS == "darwin" {
			ui.Run(ctx)
		} else {
			server.Run(ctx)
		}
	}
}
