package storage

import (
	"RobotService/internal/entities"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisCache struct {
	client *redis.Client
}

func NewClient(addr string) *RedisCache {
	rdb := redis.NewClient(
		&redis.Options{
			Addr: addr,
		})
	return &RedisCache{client: rdb}
}

func (r *RedisCache) SetRobotData(key string, robotdata entities.Robot, ttl time.Duration) error {
	pref := "robots:"
	data, err := json.Marshal(robotdata)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, pref+key, data, ttl).Err()
}

func (r *RedisCache) GetRobotData(key string) (*entities.Robot, error) {
	pref := "robots:"
	data, err := r.client.Get(ctx, pref+key).Result()
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

func (r *RedisCache) DeleteRobotData(key string) error {
	pref := "robots:"
	return r.client.Del(ctx, pref+key).Err()
}
