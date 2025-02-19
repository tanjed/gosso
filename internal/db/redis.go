package db

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"reflect"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tanjed/go-sso/internal/config"
)

type Redis struct {
	Conn *redis.Client
}

func InitRedis() *Redis {
	var rdb Redis
	rdb.Conn = redis.NewClient(&redis.Options{
        Addr:    config.AppConfig.REDIS_HOST+":"+strconv.Itoa(config.AppConfig.REDIS_PORT),
        Password: "",
        DB:       0,
    })

	return &rdb
}

func (rdb Redis) Close() {
	rdb.Conn.Close()
}


func RedisGetToStruct(key string, model interface{}) (error) {

	if !validateStruct(model) {
		slog.Error("Invalid struct type given")
		return errors.New("invalid struct type given")
	}

	rdb := InitRedis()
	defer rdb.Close()

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
	defer rdb.Close()

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