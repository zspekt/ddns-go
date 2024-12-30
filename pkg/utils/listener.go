package utils

import (
	"context"
	"log/slog"
	"net"
)

func AcceptWithCtx(l net.Listener, ctx context.Context) (net.Conn, error) {
	slog.Info("AcceptWithCtx(): called...")
	var (
		conn net.Conn
		err  error
	)

	ch := make(chan struct{}, 1)
	go func() {
		conn, err = l.Accept()
		ch <- struct{}{}
	}()
	for {
		select {
		case <-ch:
			if err != nil {
				return nil, err
			}
			return conn, nil
		case <-ctx.Done():
			slog.Info("AcceptWithCtx(): caught cancel signal while waiting for conn")
			return nil, ShutdownErr
		}
	}
}
