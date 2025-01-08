package listener

import (
	"context"
	"log"
	"net"

	"github.com/zspekt/ddns-go/pkg/utils"
	"github.com/zspekt/ddns-go/src/ip"
)

func Listen(l net.Listener, ch chan string, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Fatal("die") // TODO: implement shutdown logic
		default:
			conn, err := utils.AcceptWithCtx(l, ctx)
			if err != nil {
				if err == utils.ShutdownErr {
					log.Fatal("die") // TODO: implement shutdown logic
				}
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
