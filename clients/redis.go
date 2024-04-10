package clients

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConf struct {
	Addr         string `json:"addr" mapstructure:"addr"`
	Password     string `json:"password" mapstructure:"password"`
	DB           int    `json:"db" mapstructure:"db"`
	DialTimeout  int    `json:"dialTimeout" mapstructure:"dialTimeout"`
	WriteTimeout int    `json:"writeTimeout" mapstructure:"writeTimeout"`
	ReadTimeout  int    `json:"readTimeout" mapstructure:"readTimeout"`
	PoolTimeout  int    `json:"poolTimeout" mapstructure:"poolTimeout"`
}

type RedisClient struct {
	ctx    context.Context
	Client *redis.Client
}

func NewRedisClient(ctx context.Context, conf RedisConf) (*RedisClient, error) {
	var (
		rdb *redis.Client
		err error
	)
	rdb = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password, // 没有密码，默认值
		DB:       conf.DB,       // 默认DB 0

		PoolSize:     4 * runtime.NumCPU(),
		MinIdleConns: 2 * runtime.NumCPU(),

		DialTimeout:  time.Duration(conf.DialTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Second,
		PoolTimeout:  time.Duration(conf.PoolTimeout) * time.Second,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Connect to redis clients failed, err: %v", err))
	}
	return &RedisClient{ctx: ctx, Client: rdb}, err
}

func (rd *RedisClient) Set(key string, val interface{}, expire int) error {
	return rd.Client.Set(rd.ctx, key, val, time.Duration(expire)*time.Second).Err()
}

func (rd *RedisClient) Get(key string) (string, error) {
	return rd.Client.Get(rd.ctx, key).Result()
}

func (rd *RedisClient) Del(key string) (int64, error) {
	return rd.DelAll([]string{key}...)
}

func (rd *RedisClient) DelAll(keys ...string) (int64, error) {
	return rd.Client.Del(rd.ctx, keys...).Result()
}

func (rd *RedisClient) HSet(key string, values ...interface{}) error {
	return rd.Client.HSet(rd.ctx, key, values...).Err()
}

func (rd *RedisClient) HGet(key, field string) (string, error) {
	return rd.Client.HGet(rd.ctx, key, field).Result()
}

func (rd *RedisClient) HDel(key, field string) (int64, error) {
	return rd.HDelAll(key, []string{field}...)
}

func (rd *RedisClient) HDelAll(key string, fields ...string) (int64, error) {
	return rd.Client.HDel(rd.ctx, key, fields...).Result()
}
