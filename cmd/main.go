package main

import (
	"context"
	"fmt"
	"github.com/JMURv/seo-svc/internal/cache/redis"
	ctrl "github.com/JMURv/seo-svc/internal/controller"
	"github.com/JMURv/seo-svc/internal/discovery"

	//handler "github.com/JMURv/seo-svc/internal/handler/http"
	handler "github.com/JMURv/seo-svc/internal/handler/grpc"
	tracing "github.com/JMURv/seo-svc/internal/metrics/jaeger"
	metrics "github.com/JMURv/seo-svc/internal/metrics/prometheus"
	"go.uber.org/zap"
	//mem "github.com/JMURv/par-pro/internal/repository/memory"
	db "github.com/JMURv/seo-svc/internal/repository/db"
	cfg "github.com/JMURv/seo-svc/pkg/config"
	"os"
	"os/signal"
	"syscall"
)

const configPath = "local.config.yaml"

func mustRegisterLogger(mode string) {
	switch mode {
	case "prod":
		zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	case "dev":
		zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Panic("panic occurred", zap.Any("error", err))
			os.Exit(1)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	conf := cfg.MustLoad(configPath)
	mustRegisterLogger(conf.Server.Mode)

	// Start metrics and tracing
	go metrics.New(conf.Server.Port + 5).Start(ctx)
	go tracing.Start(ctx, conf.ServiceName, conf.Jaeger)

	dscvry := discovery.New(
		conf.SrvDiscovery.URL,
		conf.ServiceName,
		fmt.Sprintf("%v://%v:%v", conf.Server.Scheme, conf.Server.Domain, conf.Server.Port),
	)

	if err := dscvry.Register(); err != nil {
		zap.L().Warn("Error registering service", zap.Error(err))
	}

	// Setting up main app
	cache := redis.New(conf.Redis)
	repo := db.New(conf.DB)

	svc := ctrl.New(repo, cache)
	h := handler.New(svc)

	// Graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-c

		zap.L().Info("Shutting down gracefully...")

		cancel()
		cache.Close()
		if err := dscvry.Deregister(); err != nil {
			zap.L().Warn("Error deregistering service", zap.Error(err))
		}
		if err := h.Close(); err != nil {
			zap.L().Debug("Error closing handler", zap.Error(err))
		}

		os.Exit(0)
	}()

	// Start service
	zap.L().Info(
		fmt.Sprintf("Starting server on %v://%v:%v", conf.Server.Scheme, conf.Server.Domain, conf.Server.Port),
	)
	h.Start(conf.Server.Port)
}
