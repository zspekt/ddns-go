package listener

import (
	"context"
	"log"
	"log/slog"
	"net"

	"github.com/zspekt/ddns-go/pkg/utils"
	"github.com/zspekt/ddns-go/src/ip"
)

func Listen(l net.Listener, ch chan string, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Fatal("Listen(): AcceptWithCtx() received shutdown signal")
		default:
			conn, err := utils.AcceptWithCtx(l, ctx)
			if err != nil {
				if err == utils.ShutdownErr {
					log.Fatal("Listen(): AcceptWithCtx() received shutdown signal")
				}
				slog.Error("Listen(): unexpected error. continuing loop...", "error", err)
				continue
			}
			handleConn(conn, ch)
		}
	}
}

func handleConn(conn net.Conn, ch chan string) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	check(err) // TODO: handle error
	buf = buf[:n]

	ip, err := ip.ParseIP(buf)
	check(err) // TODO: handle error

	ch <- ip
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
