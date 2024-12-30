package main

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/zspekt/ddns-go/pkg/utils"
	"github.com/zspekt/ddns-go/src/dns"
)

/*
	pseudocode

create listener

	for {
	  listen on addr:port

	  if address not stored {
	    send put ddns
	    store it
	    back to loop (waiting for conn)
	  }

	  if address stored {
	    compare
	    if new addr is equal to old {
	      back to loop
	    }
	    if new addr is *NOT* equal to old {
	      send put ddns
	      store it
	      back to loop (waiting for conn)
	    }
	  }
	}
*/

func main() {
	// l, err := net.Listen("tcp", "192.168.1.162:33219")
	// must(err)
	// ctx, cancel := context.WithCancel(context.Background())
	//
	// sigs := make(chan os.Signal, 1)
	// go shutdown(sigs, cancel)
	//
	// file, err := os.OpenFile("ip.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	// print(file, err)
	//
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		log.Fatal("die") // TODO: implement shutdown logic
	// 	default:
	// 		conn, err := AcceptWithCtx(l, ctx)
	// 		if err != nil {
	// 			if err == shutdownErr {
	// 				// TODO: addiational shutdown logic
	// 			}
	// 		}
	//
	// 		reader := bufio.NewReader(conn)
	// 		b, err := reader.ReadBytes('\n')
	// 		print(b)
	// 		must(
	// 			err,
	// 		) // we should always reach the delimiter, so no reason why we should ever get io.EOF
	// 	}
	// }
	// ctx := context.Background()
	utils.Must(godotenv.Load(".env"))

	token := os.Getenv("CLOUDFLARE_API_TOKEN") // TODO: config funcs to set this up

	dns.UpdateDnsRecord("177.177.177.177", token)
}

/*


























































 */
