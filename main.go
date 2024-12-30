package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kr/pretty"
)

var shutdownErr error = errors.New("got shutdown signal")

const (
	FILE_EXISTS  int = 1
	FILE_CREATED int = 2
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func shutdown(sigs chan os.Signal, cancel context.CancelFunc) {
	slog.Info("shutdown(): starting routine...")

	// we register the channel so it will get these sigs
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	s := <-sigs
	slog.Info(
		fmt.Sprintf(
			"shutdown routine caught %v sig. cancelling...",
			s.String(),
		),
	)
	cancel()
}

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
			return nil, shutdownErr
		}
	}
}

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
	l, err := net.Listen("tcp", "192.168.1.162:33219")
	must(err)
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	go shutdown(sigs, cancel)

	file, err := os.OpenFile("ip.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	print(file, err)

	for {
		select {
		case <-ctx.Done():
			log.Fatal("die") // TODO: implement shutdown logic
		default:
			conn, err := AcceptWithCtx(l, ctx)
			if err != nil {
				if err == shutdownErr {
					// TODO: addiational shutdown logic
				}
			}

			reader := bufio.NewReader(conn)
			b, err := reader.ReadBytes('\n')
			print(b)
			must(
				err,
			) // we should always reach the delimiter, so no reason why we should ever get io.EOF
		}
	}
}

// decodes reader with JSON-formatted content r into the STRUCT pointed at by st
func DecodeJson[T any](r io.Reader, st *T) error {
	decoder := json.NewDecoder(r)

	// ipToChangeTo := "190.221.132.80"

	err := decoder.Decode(st)
	if err != nil {
		slog.Error("error decoding json", "error", err)
		return err
	}
	return nil
}

func clean(ctx context.Context) {
	client := http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest(
		"GET",
		"https://api.cloudflare.com/client/v4/zones/",
		nil,
	)
	req.Header.Set("Authorization", "Bearer XXXXXXX")

	var getZonesIDs *dnsGetZonesIDs = new(dnsGetZonesIDs)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	err = DecodeJson(res.Body, getZonesIDs)
	if err != nil {
		log.Fatal("getZonesIDs", err)
	}
	zoneID := getZonesIDs.Result[0].ID

	/////////////////////////////////////////////////////////////////////////////

	req, err = http.NewRequest(
		"GET",
		"https://api.cloudflare.com/client/v4/zones/"+zoneID+"/dns_records/",
		nil,
	)
	req.Header.Set("Authorization", "Bearer XXXXXXX")

	res, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var getDnsRecords *dnsGetRecords = new(dnsGetRecords)
	err = DecodeJson(res.Body, getDnsRecords)
	if err != nil {
		log.Fatal("getDnsRecords", err)
	}

	for _, dnsRecord := range getDnsRecords.Result {
		r, err := dnsPutRequest(
			dnsRecord,
			"190.221.132.80",
			"XXXXXXX",
		)
		if err != nil {
			log.Fatal(err)
		}
		respp := new(dnsPutResponse)
		err = DecodeJson(io.Reader(r.Body), respp)
		if err != nil {
			log.Fatal("error decoding final response")
		}
		fmt.Println(pretty.Println(respp))
	}
}

func dnsPutRequest(rec dnsRecord, ip, token string) (*http.Response, error) {
	d := ddnsPut{
		Comment:  "go ddns lol",
		Name:     rec.Name,
		Proxied:  rec.Proxied,
		Settings: rec.Settings,
		Tags:     rec.Tags,
		TTL:      rec.TTL,
		Content:  ip,
		Type:     rec.Type,
	}
	jsonPayload, err := json.Marshal(d)
	r := bytes.NewReader(jsonPayload)

	url := "https://api.cloudflare.com/client/v4/zones/" + rec.ZoneID + "/dns_records/" + rec.ID
	req, err := http.NewRequest(
		"PUT",
		url,
		r,
	)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "appliaction/json; charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{
		Timeout: 15 * time.Second,
	}
	return client.Do(req)
}

// high level routine. is responsible for receiving IP string, storing, comparing
// and updating if necessary (through PUT request)
func putGoRoutine(ctx context.Context, c <-chan string, token string)

// compares given IP to the stored value, updating it if necessary. returns true
// only if given IP matches the value, if no value was stored (file didn't exist),
// or it differs from the one passed in, it returns false.
func compareAndStoreIP() bool

// makes the put request to cloudflare's api to update all A records to point
// to the ip that's been passed in. returns any error found throughout the execution
// of the function, or any error from cloudflare's response
func putRequest(ip, token string) error { // FIX: this function is too long
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	var (
		baseURL       string          = "https://api.cloudflare.com/client/v4"
		getZonesIDs   *dnsGetZonesIDs = new(dnsGetZonesIDs)
		getDnsRecords *dnsGetRecords  = new(dnsGetRecords)
	)

	req, err := newRequestWithToken("GET", baseURL+"/zones/", token, nil)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	err = DecodeJson(res.Body, getZonesIDs)
	if err != nil {
		log.Fatal("getZonesIDs", err)
	}
	zoneID := getZonesIDs.Result[0].ID

	/////////////////////////////////////////////////////////////////////////////

	req, err = newRequestWithToken("GET", baseURL+"/zones/"+zoneID+"/dns_records/", token, nil)

	res, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	err = DecodeJson(res.Body, getDnsRecords)
	if err != nil {
		log.Fatal("getDnsRecords", err)
	}

	for _, dnsRecord := range getDnsRecords.Result {
		r, err := dnsPutRequest(
			dnsRecord,
			"190.221.132.80",
			"XXXXXXX",
		)
		// FIX:   only try to replace "A" records content.
		//        don't currently have any others, but it's stinky code.
		//
		if err != nil {
			log.Fatal(err)
		}
		respp := new(dnsPutResponse)
		err = DecodeJson(io.Reader(r.Body), respp)
		if err != nil {
			log.Fatal("error decoding final response")
		}
	}
	return nil
}

func newRequestWithToken(method, url, token string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", "Bearer "+token)
	return r, nil
}

/*


























































 */
