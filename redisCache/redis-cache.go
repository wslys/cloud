package redisCache

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"time"
)

type RedisCLI struct {}

var rsv *redis.Pool
var ctx context.Context
func init() {
	rsv = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	ctx = context.TODO()
}

func (cli *RedisCLI) Get(key string, result interface{}) error {
	client := rsv.Get()

	if client == nil {
		return errors.InternalServerError("Get Redis", "redis connection is nil")
	}
	defer client.Close()

	valueBytes, err := redis.Bytes(client.Do("GET", key))
	if err != nil {
		return errors.InternalServerError("Get Redis", err.Error())
	}
	err = json.Unmarshal(valueBytes, result)
	if err != nil {
		return errors.InternalServerError("Get Redis", err.Error())
	}
	return nil
}

func (cli *RedisCLI) Set(key string, data interface{}) error {
	client := rsv.Get()
	if client == nil {
		return errors.InternalServerError("Redis Set", "redis connection is nil")
	}
	defer client.Close()
	value, err := json.Marshal(data)

	if err != nil {
		return errors.InternalServerError("Redis Set", err.Error())
	}
	_, err = client.Do("SET", key, value)
	if err != nil {
		return errors.InternalServerError("Redis Set", err.Error())
	}
	return nil
}

func (cli *RedisCLI) Update(key string, data interface{}) error {
	client := rsv.Get()
	if client == nil {
		return errors.InternalServerError("Redis Update", "redis connection is nil")
	}
	defer client.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return errors.InternalServerError("Redis Update", err.Error())
	}
	_, err = client.Do("DEL", key)
	if err != nil {
		return errors.InternalServerError("Redis Update", err.Error())
	}

	_, err = client.Do("SET", key, value)
	if err != nil {
		return errors.InternalServerError("Redis Update", err.Error())
	}
	return nil
}

func (cli *RedisCLI) Delete(key string) error {
	client := rsv.Get()
	if client == nil {
		return errors.InternalServerError("Redis Delete", "redis connection is nil")
	}
	defer client.Close()
	_, err := client.Do("DEL", key)
	if err != nil {
		return errors.InternalServerError("Redis Delete", err.Error())
	}
	return nil
}

// TODO
func (cli *RedisCLI) Clear() error {
	return nil
}
