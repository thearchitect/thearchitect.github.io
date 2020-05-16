package pty

import (
	"context"
	"io"
	"log"
	"time"
)

func (p *PTY) Attach(ctx context.Context, remote io.ReadWriter) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	done := func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
			return false
		}
	}

	go func() {
		defer cancel()

		var buf [32 * 1024]byte
		for {
			if done() {
				return
			}

			n, err := p.Master.Read(buf[:])
			if err != nil {
				log.Println(err)
				time.Sleep(500 * time.Millisecond)
				continue
			}

			b := buf[:n]

			log.Println("output:", len(b))
			if _, err := remote.Write(b); err != nil {
				log.Println(err)
				time.Sleep(500 * time.Millisecond)
				return
			}
		}
	}()

	go func() {
		defer cancel()

		var buf [32 * 1024]byte
		for {
			if done() {
				return
			}

			n, err := remote.Read(buf[:])
			if err != nil {
				log.Println(err)
				time.Sleep(500 * time.Millisecond)
				return
			}

			b := buf[:n]

			log.Println("input:", string(b))
			if _, err := p.Master.Write(b); err != nil {
				log.Println(err)
				time.Sleep(500 * time.Millisecond)
				continue
			}

		}
	}()

	<-ctx.Done()
}
