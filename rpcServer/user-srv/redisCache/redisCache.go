package redisCache

import (
	"beaconCloud/rpcServer/user-srv/proto/loginCache"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"time"
)

type RedisCLI struct {
	rs *redis.Pool
}

var rsv *redis.Pool
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
}


func (cli *RedisCLI) GetLogin(ctx context.Context, req *loginCache.GetLoginRequest, rsp *loginCache.GetLoginResponse) error {
	client := rsv.Get()
	if client == nil {
		return errors.InternalServerError("go.micro.srv.user.GetLogin", "redis connection is nil")
	}
	defer client.Close()
	valueBytes, err := redis.Bytes(client.Do("GET", req.Login.UserId))
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.GetLogin", err.Error())
	}
	login := &loginCache.Login{}
	err = json.Unmarshal(valueBytes, login)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.GetLogin", err.Error())
	}
	rsp.Login = login
	return nil
}

func (cli *RedisCLI) AddLogin(ctx context.Context, req *loginCache.AddLoginRequest, rsp *loginCache.AddLoginResponse) error {
	client := rsv.Get()
	if client == nil {
		return errors.InternalServerError("go.micro.srv.user.AddLogin", "redis connection is nil")
	}
	defer client.Close()
	value, err := json.Marshal(req.Login)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.AddLogin", err.Error())
	}
	_, err = client.Do("SET", req.Login.UserId, value)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.AddLogin", err.Error())
	}
	return nil
}

func (cli *RedisCLI) ActivateLogin(ctx context.Context, req *loginCache.ActivateLoginRequest, rsp *loginCache.ActivateLoginResponse) error {
	client := rsv.Get()
	if client == nil {
		return errors.InternalServerError("go.micro.srv.user.GetLogin", "redis connection is nil")
	}
	defer client.Close()
	req.Login.Created = time.Now().Unix()
	req.Login.Validity = 24 * 7 * 3600
	req.Login.ValidityTime = time.Now().Unix() + 24*7*3600
	value, err := json.Marshal(req.Login)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.ActivateLogin", err.Error())
	}
	_, err = client.Do("SET", req.Token, value)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.ActivateLogin", err.Error())
	}
	return nil
}

func (cli *RedisCLI) DeleteLogin(ctx context.Context, req *loginCache.DeleteLoginRequest, rsp *loginCache.DeleteLoginResponse) error {
	client := rsv.Get()
	if client == nil {
		return errors.InternalServerError("go.micro.srv.user.DeleteLogin", "redis connection is nil")
	}
	defer client.Close()
	_, err := client.Do("DEL", req.Login.UserId)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.DeleteLogin", err.Error())
	}
	return nil
}

func (cli *RedisCLI) Clear(ctx context.Context, req *loginCache.ClearRequest, rsp *loginCache.ClearResponse) error {

	return nil
}
