package setup

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/zspekt/ddns-go/pkg/utils"
	"github.com/zspekt/ddns-go/src/dns"
)

type Cfg struct {
	Listener net.Listener
	Api      dns.Api
	Filename string
	Ctx      context.Context
	Cancel   context.CancelFunc
	Sigs     chan os.Signal
	Ch       chan string
}

/*
	NOTE:

envs:

	CLOUDFLARE_API_TOKEN   -> api token
	LOG_LEVEL  (def INFO)  -> log level
	ADDR  (def localhost)  -> address to listen on
	PORT       (def 8080)  -> port to listen on
	FILENAME (def ip.txt)  -> file to store the IP value in
*/
func Config() *Cfg {
	logger()

	token := token()
	baseURL := "https://api.cloudflare.com/client/v4"

	api := &dns.CloudFlareAPI{
		Token:   token,
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 6 * time.Second,
		},
	}

	addr := utils.EnvOrDefault("ADDR", "localhost")
	port := utils.EnvOrDefault("PORT", "8080")
	l, err := net.Listen("tcp", addr+":"+port)
	if err != nil {
		log.Fatal("error creating listener", err)
	}

	filename := utils.EnvOrDefault("FILENAME", "ip.txt")
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	ch := make(chan string, 1)

	return &Cfg{
		Listener: l,
		Api:      api,
		Filename: filename,
		Ctx:      ctx,
		Cancel:   cancel,
		Sigs:     sigs,
		Ch:       ch,
	}
}

func logger() {
	// https://stackoverflow.com/a/76970969
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)
	slog.SetDefault(
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})),
	)

	l, ok := os.LookupEnv("LOG_LEVEL")
	if !ok || l == "" {
		slog.Info("no log level provided. defaulting to INFO...")
		return
	}
	logLevels := map[string]slog.Level{"DEBUG": -4, "INFO": 0, "WARN": 4, "ERROR": 8}

	v, ok := logLevels[l]
	if !ok {
		slog.Info("invalid log level provided. defaulting to INFO...", "logLevel", l)
		return
	}
	lvl.Set(v)
	slog.Info("log level set", "logLevel", l)
}

func token() string {
	token, ok := os.LookupEnv("CLOUDFLARE_API_TOKEN")
	if ok && token != "" {
		return strings.TrimSpace(token)
	}

	slog.Info("no token env var found")

	b, err := os.ReadFile("/run/secrets/CLOUDFLARE_API_TOKEN")
	if err != nil {
		log.Fatal("no token found in .env or in /run/secrets")
	}

	return strings.TrimSpace(string(b))
}
