package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mashmorsik/banners-service/config"
	"github.com/mashmorsik/banners-service/infrastructure/data"
	"github.com/mashmorsik/banners-service/infrastructure/data/cache"
	"github.com/mashmorsik/banners-service/infrastructure/server"
	"github.com/mashmorsik/banners-service/internal/banner"
	"github.com/mashmorsik/banners-service/pkg/token"
	"github.com/mashmorsik/banners-service/repository"
	"github.com/mashmorsik/logger"
)

func main() {
	logger.BuildLogger(nil)

	conf, err := config.LoadConfig()
	if err != nil {
		logger.Errf("Error loading config: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGKILL)

	go func() {
		<-sigCh
		logger.Infof("context done")
		cancel()
	}()

	conn := data.MustConnectPostgres(ctx, conf)
	data.MustMigrate(conn)

	dat := data.NewData(ctx, conn)

	bannerCache := cache.NewBannerCache(ctx, conf.Cache.EvictionWorkerDuration, conf)

	bannerRepo := repository.NewBannerRepo(ctx, dat)
	bb := banner.NewBanner(ctx, bannerRepo, conf, &bannerCache)

	token.NewTokenManager(conf.Auth.TokenSecret)

	httpServer := server.NewServer(conf, *bb)
	if err = httpServer.StartServer(ctx); err != nil {
		logger.Warn(err.Error())
	}
}
