package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"example.com/order-service/httptransport"
	"example.com/order-service/implementation"
	"example.com/order-service/repository"
	"example.com/order-service/transport"
)

func main() {
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	logger = level.NewFilter(logger, level.AllowDebug())
	logger = log.With(
		logger,
		"svc", "order",
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)

	repo := repository.NewMemoryRepository()
	svc := implementation.NewService(repo, logger)
	endpoints := transport.MakeEndpoints(svc)
	handler := httptransport.NewService(endpoints, logger)

	server := &http.Server{
		Addr:    *httpAddr,
		Handler: handler,
	}

	serverErrors := make(chan error, 1)

	go func() {
		level.Info(logger).Log("msg", "service started", "addr", *httpAddr)
		serverErrors <- server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(stop)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			level.Error(logger).Log("msg", "server error", "err", err)
		}
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			level.Error(logger).Log("msg", "shutdown error", "err", err)
		}
	}
}
