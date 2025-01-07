package utils

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var ShutdownErr error = errors.New("received shutdown signal")

func Shutdown(sigs chan os.Signal, cancel context.CancelFunc) {
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
