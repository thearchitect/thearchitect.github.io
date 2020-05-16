package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"

	"github.com/thearchitect/thearchitect.github.io/internal/pty"
	"github.com/thearchitect/thearchitect.github.io/server/resources"
)

func ServeContainerHTTP() http.HandlerFunc {
	var upgrader = websocket.Upgrader{
		EnableCompression: true,
		CheckOrigin: func(q *http.Request) bool {
			return true
		},
	}

	var dockerContext = resources.DockerContext()

	return func(w http.ResponseWriter, q *http.Request) {
		conn, err := upgrader.Upgrade(w, q, http.Header{
			"OK": {"true"},
		})
		if err != nil {
			panic(err)
		}

		ServeContainer(q.Context(), dockerContext, conn)
	}
}

func ServeContainer(ctx context.Context, dockerContext func() io.Reader, conn *websocket.Conn) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pty, err := pty.New()
	if err != nil {
		panic(err)
	}
	defer pty.Close()

	stdio := io.ReadWriter(&WSIO{
		conn:   conn,
		prefix: []byte("\x1b]"),
		catcher: func(data []byte) bool {
			data = data[2:]

			var cmd = struct {
				W int `json:"w"`
				H int `json:"h"`
			}{}
			if err := json.Unmarshal(data, &cmd); err != nil {
				panic(err)
			}

			log.Println("resize", cmd)

			if err := pty.Resize(cmd.W, cmd.H); err != nil {
				panic(err)
			}

			return false
		},
	})

	go func() {
		defer cancel()

		pty.Attach(ctx, stdio)
	}()

	if cmd := pty.Command(ctx, "docker", "build", "--build-arg", "USERNAME=neo", "-t", "zion", "-"); cmd != nil {
		cmd.Stdin = dockerContext()

		Banner(stdio, "Building container")
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

	// TODO https://docs.docker.com/config/containers/resource_constraints/
	Banner(stdio, "Starting container")
	if err := pty.Command(ctx,
		"docker", "run",
		"--rm", "-it",
		"--network=none",
		"--memory=32m",
		"--memory-swap=64m",
		"--oom-kill-disable",
		"--cpus=2", "--cpu-shares=512",
		"-h=zion",
		"zion",
	).Run(); err != nil {
		panic(err)
	}

}

////////////////////////////////////////////////////////////////
//// Banner
////

func Banner(w io.Writer, args ...interface{}) {

	line := fmt.Sprintf("\r\n%s\r\n", strings.Repeat("•", 64))

	args = append([]interface{}{line}, args...)
	args = append(args, line)

	if _, err := color.New(color.FgGreen).Fprint(w, args...); err != nil {
		panic(err)
	}
}

////////////////////////////////////////////////////////////////
//// WSIO
////

var _ io.ReadWriter = new(WSIO)

type WSIO struct {
	conn    *websocket.Conn
	prefix  []byte
	catcher func([]byte) bool
}

func (io *WSIO) Read(p []byte) (n int, err error) {
	for {
		_, data, err := io.conn.ReadMessage()
		if err != nil {
			return 0, err
		}

		if len(io.prefix) > 0 && io.catcher != nil && bytes.HasPrefix(data, io.prefix) {
			if io.catcher(data) {
				n = copy(p, data)
				return n, nil
			}
		} else {
			n = copy(p, data)
			return n, nil
		}
	}
}

func (io *WSIO) Write(p []byte) (n int, err error) {
	err = io.conn.WriteMessage(websocket.BinaryMessage, p)
	n = len(p)
	return
}
