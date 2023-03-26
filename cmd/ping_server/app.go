package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go_ping_kube/internal/config"
	"go_ping_kube/internal/delivery/rest"
	"go_ping_kube/internal/domain/services"
	"go_ping_kube/internal/infrastructure/repository"
)

func runPingApp() error {
	ctx := context.TODO()

	logrus.Info("app running...")

	appConfig, err := config.NewConfig("configs/dev_config.yaml")
	if err != nil {
		return err
	}
	logrus.Info("init config >>> OK")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(":%s", appConfig.Redis.RedisAddr),
		Password: appConfig.Redis.RedisPass,
		DB:       appConfig.Redis.RedisDB,
	})
	logrus.Info("redis client >>> OK")
	pingRepository, err := repository.NewProxyRedisStore(ctx, redisClient)
	if err != nil {
		return err
	}

	pingService := services.NewPingService(pingRepository)

	pingHandler := rest.NewPingHandler(pingService)

	pingServer := rest.NewPingServer(pingHandler)

	return pingServer.Start(appConfig.AppPort)
}
