package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/1ocknight/mess/shared/logger"
	"github.com/1ocknight/mess/shared/verify"
	"github.com/1ocknight/mess/websocket/config"
	"github.com/1ocknight/mess/websocket/internal/ctxkey"
	"github.com/1ocknight/mess/websocket/internal/hub/general"
	"github.com/1ocknight/mess/websocket/internal/loglables"
	"github.com/1ocknight/mess/websocket/internal/transport"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lg := logger.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	lg = lg.With(loglables.Service, "websocket_microservice")

	ctx = ctxkey.WithLogger(ctx, lg)

	cfg, err := config.LoadConfig()
	if err != nil {
		lg.Error(fmt.Errorf("load config: %w", err))
		return
	}

	verLg := lg.With(loglables.Layer, "verify")
	ver, err := verify.New(cfg.Verify, verLg)
	if err != nil {
		lg.Error(fmt.Errorf("verify new: %w", err))
		return
	}

	ghLg := lg.With(loglables.Layer, "general_hub")
	ghub := general.NewHub(ctx, ghLg, cfg.GeneralHub)
	ghub.Start(ctx)

	handler := transport.NewHandler(cfg.Handler, ghub)

	serverLg := lg.With(loglables.Layer, "server")
	server := transport.NewServer(cfg.HTTP, ver, handler, serverLg)
	go func() {
		if err := server.Run(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			lg.Error(fmt.Errorf("server run: %w", err))
			return
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	lg.Info("start graceful shutdown")

	err = server.Stop(ctx)
	if err != nil {
		lg.Error(fmt.Errorf("server stop: %w", err))
	}
	lg.Info("server is stop")

	cancel()
	lg.Info("successful stop")
}
