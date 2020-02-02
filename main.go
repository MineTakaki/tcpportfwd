package main

import (
	"context"
	"io"
	"net"
	"os"
	"time"

	"go.uber.org/zap"
)

type (
	contextKey int
)

const (
	keyOfLog contextKey = iota
)

func main() {
	os.Exit(mainProc())
}

func mainProc() int {
	l := zap.L()
	defer l.Sync()

	ctx, fnCancel := context.WithCancel(context.Background())
	defer fnCancel()

	ctx = context.WithValue(ctx, keyOfLog, l)

	lcfg := net.ListenConfig{}
	ln, err := lcfg.Listen(ctx, "tcp", "localhost:1522")
	if err != nil {
		return 9
	}

	for {
		select {
		default:
		case <- ctx.Done():
			break
		}

		conn, err := ln.Accept()
		if err != nil {
			return 8
		}

		go handle(ctx, conn, "localhost:1521")
	}

	return 0
}

func forward(src, dst net.Conn) {
	defer src.Close()
	defer dst.Close()
	io.Copy(src, dst)
}

func handle(ctx context.Context, c net.Conn, address string) {
	l, _ := ctx.Value(keyOfLog).(*zap.Logger)
	l.Info("connection from", zap.String("address", c.RemoteAddr().String()))

	d := net.Dialer{
		Timeout: 30 * time.Second,
	}
	rmt, err := d.DialContext(ctx, "tcp", address)
	if err != nil {
	}

	l.Info("connection to", zap.String("address", address))

	go forward(c, rmt)
	go forward(rmt, c)
}
