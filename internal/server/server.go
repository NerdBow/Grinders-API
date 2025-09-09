package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()

	server := http.Server{
		Addr:              os.Getenv("ADDRESS"),
		Handler:           mux,
		ReadHeaderTimeout: time.Second,
	}

	slog.Info("Starting server.", slog.String("Address", server.Addr))

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			slog.Error("Unable to start server.", slog.String("error", err.Error()))
			stop()
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(context.Background()); err != nil {
		slog.Error("Unable to gracefully shutdown server.", slog.String("error", err.Error()))
	}
}
