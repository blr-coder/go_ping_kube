package rest

import (
	"context"
	"fmt"
	"github.com/gavv/httpexpect/v2"
	"github.com/redis/go-redis/v9"
	"go_ping_kube/internal/domain/services"
	"go_ping_kube/internal/infrastructure/repository"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingServer_handleHttp(t *testing.T) {

	ctx := context.TODO()

	/*appConfig, err := config.NewConfig("/configs/dev_config.yaml")
	if err != nil {
		return
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(":%s", appConfig.Redis.RedisAddr),
		Password: appConfig.Redis.RedisPass,
		DB:       appConfig.Redis.RedisDB,
	})*/

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(":%s", "6379"),
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
		DB:       0,
	})

	pingRepository, err := repository.NewProxyRedisStore(ctx, redisClient)
	if err != nil {
		return
	}

	pingService := services.NewPingService(pingRepository)
	pingHandler := NewPingHandler(pingService)
	pingServer := NewPingServer(pingHandler)

	server := httptest.NewServer(pingServer.handleHttp())
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.GET("/api/v1/get").
		WithQuery("uuid", "990d45f5-d7f2-41ff-8597-2651d8cff2d1").
		Expect().
		Status(http.StatusOK).
		JSON().
		NotNull()

	e.GET("/api/v1/get").
		WithQuery("uuid", "990d45f5-d7f2-41ff-8597-2651d8cff2d2").
		Expect().
		Status(http.StatusNotFound).
		JSON()

	e.POST("/api/v1/get").
		WithQuery("uuid", "990d45f5-d7f2-41ff-8597-2651d8cff2d1").
		Expect().
		Status(http.StatusMethodNotAllowed).
		JSON()
}
