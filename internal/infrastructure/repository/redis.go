package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go_ping_kube/internal/domain/models"
)

type PingRedisStore struct {
	Client *redis.Client
}

func NewProxyRedisStore(ctx context.Context, client *redis.Client) (*PingRedisStore, error) {
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	logrus.Info("redis PING >>> OK")

	return &PingRedisStore{
		Client: client,
	}, nil
}

func (s *PingRedisStore) Save(ctx context.Context, ping *models.PingData) error {
	fmt.Println("REDIS Save")

	err := s.Client.HSet(ctx, "ping_hash", ping.Id.String(), ping).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *PingRedisStore) Get(ctx context.Context, uuid uuid.UUID) (*models.PingData, error) {
	fmt.Println("REDIS Get")

	pingData := &models.PingData{}

	data := s.Client.HGet(ctx, "ping_hash", uuid.String())
	if data.Err() != nil {
		return nil, data.Err()
	}

	err := json.Unmarshal([]byte(data.Val()), &pingData)
	if err != nil {
		return nil, err
	}

	/*pingData := &models.PingData{}
	if err := s.Client.HGet(ctx, "ping_hash", uuid.String()).Scan(pingData); err != nil {
		panic(err)
	}*/

	return pingData, nil
}

func (s *PingRedisStore) All(ctx context.Context) ([]*models.PingData, error) {
	list := s.Client.HGetAll(ctx, "ping_hash")
	if list.Err() != nil {
		return nil, list.Err()
	}

	var dataList []*models.PingData
	for _, ping := range list.Val() {

		pingData := &models.PingData{}

		err := json.Unmarshal([]byte(ping), &pingData)
		if err != nil {
			return nil, err
		}

		dataList = append(dataList, pingData)
	}

	return dataList, nil
}
