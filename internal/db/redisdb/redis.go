package redisdb

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tanjed/go-sso/internal/config"
)

type Redis struct {
	id int
	Conn *redis.Client
}

var conn *redis.Client

func NewRedis(c *config.Config) *Redis {
	 client := Redis{
			id: rand.Intn(101),
			Conn: redis.NewClient(&redis.Options{
				Addr:    c.REDIS_HOST+":"+strconv.Itoa(c.REDIS_PORT),
				Password: "",
				DB:       0,
				PoolSize:     500,                                     
				MinIdleConns: 50,                                           
				MaxRetries:   3,
		}),
	}

	conn = client.Conn
	return &client
}

func (c Redis) Close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}


func RedisGetToStruct(key string, model interface{}) (error) {

	if !validateStruct(model) {
		slog.Error("Invalid struct type given")
		return errors.New("invalid struct type given")
	}
	
	m, err := conn.Get(context.Background(), key).Result()

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

	err = conn.Set(context.Background(), key, jsonData, ttl).Err()
	if err != nil {
		slog.Error("Unable to set to redis", "error", err)
		return errors.New("unable to set to redis")
	}

	return nil
}


func validateStruct(s any) bool {
	return reflect.TypeOf(s).Kind() == reflect.Ptr && reflect.TypeOf(s).Elem().Kind() == reflect.Struct
}