package sorrage

import (
	"RobotService/internal/entities"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RdsCache struct {
	client *redis.Client
}

func NewClient(addr string) *RdsCache {
	rdb := redis.NewClient(
		&redis.Options{
			Addr: addr,
		})
	return &RdsCache{client: rdb}
}

func (rds *RdsCache) SetRobotData(key string, robotdata entities.Robot, ttl time.Duration) error {
	pref := "robots:"
	data, err := json.Marshal(robotdata)
	if err != nil {
		return err
	}
	return rds.client.Set(ctx, pref+key, data, ttl).Err()
}

func (rds *RdsCache) GetRobotData(key string) (*entities.Robot, error) {
	pref := "robots:"
	data, err := rds.client.Get(ctx, pref+key).Result()
	if err != nil {
		return nil, err
	}

	var robot entities.Robot
	err = json.Unmarshal([]byte(data), &robot)
	if err != nil {
		return nil, err
	}
	return &robot, nil
}

func (rds *RdsCache) DeleteRobotData(key string) error {
	pref := "robots:"
	return rds.client.Del(ctx, pref+key).Err()
}
