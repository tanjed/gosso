package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tanjed/go-sso/internal/config"
)

type Redis struct {
	Id int
	Conn *redis.Client
}

var redisInstance *Redis
var redisIntanceOnce sync.Once

func InitRedis() *Redis {
	redisIntanceOnce.Do(func ()  {
		redisInstance = &Redis{
			Id: rand.Intn(101),
			Conn : redis.NewClient(&redis.Options{
				Addr:    config.AppConfig.REDIS_HOST+":"+strconv.Itoa(config.AppConfig.REDIS_PORT),
				Password: "",
				DB:       0,
				PoolSize:     500,                                     // Maximum number of connections in the pool
				MinIdleConns: 50,                                      // Minimum idle connections                     // Idle timeout for connections
				MaxRetries:   3,                                      // Retry count for failed commands
			}),
		}
	})
	fmt.Println("REDIS CONNECTION ID:", redisInstance.Id)
	return redisInstance
}

func CloseRedis() {
	if redisInstance != nil && redisInstance.Conn != nil {
		redisInstance.Conn.Close()
	}
}


func RedisGetToStruct(key string, model interface{}) (error) {

	if !validateStruct(model) {
		slog.Error("Invalid struct type given")
		return errors.New("invalid struct type given")
	}

	rdb := InitRedis()

	m, err := rdb.Conn.Get(context.Background(), key).Result()

	if err != nil && err != redis.Nil {
		slog.Error("Unable to fetch data from redis", "error", err)
		return err
	}

	if err == redis.Nil {
		return redis.Nil
	}

	json.Unmarshal([]byte(m), &model)
	return nil

}

func RedisSetToStruct(key string, model interface{}, ttl time.Duration) error {
	if !validateStruct(model) {
		slog.Error("Invalid struct type given")
		return errors.New("invalid struct type given")
	}

	jsonData, err := json.Marshal(model)

	if err != nil {
		slog.Error("Unable to marshal data", "error", err)
		return err
	}

	rdb := InitRedis()

	err = rdb.Conn.Set(context.Background(), key, jsonData, ttl).Err()
	if err != nil {
		slog.Error("Unable to set to redis", "error", err)
		return errors.New("unable to set to redis")
	}

	return nil
}


func validateStruct(s any) bool {
	return reflect.TypeOf(s).Kind() == reflect.Ptr && reflect.TypeOf(s).Elem().Kind() == reflect.Struct
}