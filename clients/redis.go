package clients

import (
	"context"
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

func NewRedisDB(conf RedisConf) *redis.Client {
	var (
		err error
		rdb *redis.Client
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
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Errorf("connect to redis clients failed, err: %v", err))
	}
	return rdb
}

type RedisClient struct {
	ctx context.Context
	db  *redis.Client
}

func NewRedisClient(ctx context.Context, cfg RedisConf) *RedisClient {
	return &RedisClient{
		ctx: ctx,
		db:  NewRedisDB(cfg),
	}
}

func (rc *RedisClient) Set(key string, val any, expire int) error {
	return rc.db.Set(rc.ctx, key, val, time.Duration(expire)*time.Second).Err()
}

func (rc *RedisClient) Get(key string) (string, error) {
	return rc.db.Get(rc.ctx, key).Result()
}

func (rc *RedisClient) Del(key string) (int64, error) {
	return rc.DelAll([]string{key}...)
}

func (rc *RedisClient) DelAll(keys ...string) (int64, error) {
	return rc.db.Del(rc.ctx, keys...).Result()
}

func (rc *RedisClient) HSet(key string, values ...any) error {
	return rc.db.HSet(rc.ctx, key, values...).Err()
}

func (rc *RedisClient) HGet(key, field string) (string, error) {
	return rc.db.HGet(rc.ctx, key, field).Result()
}

func (rc *RedisClient) HDel(key, field string) (int64, error) {
	return rc.HDelAll(key, []string{field}...)
}

func (rc *RedisClient) HDelAll(key string, fields ...string) (int64, error) {
	return rc.db.HDel(rc.ctx, key, fields...).Result()
}
