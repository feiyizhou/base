package main

import (
	"context"
	"fmt"

	"github.com/feiyizhou/base/clients"
	msgcenter "github.com/feiyizhou/base/msgCenter"
)

func main() {
	if _, err := msgcenter.NewProducer(context.Background(), clients.NewRedisClient(
		context.Background(),
		clients.RedisConf{
			Addr:         "172.16.198.120:32763",
			Password:     "LinCloud1!2@3#4$",
			DB:           0,
			DialTimeout:  5,
			WriteTimeout: 5,
			ReadTimeout:  5,
			PoolTimeout:  5,
		}),
		"topic",
		"group",
	).Publish(msgcenter.Msg{
		MsgID:     "msgid",
		MsgType:   "msgtype",
		Source:    "source",
		Timestamp: "timestamp",
	}.ToMap()); err != nil {
		fmt.Printf("publish message failed, err: %v\n", err)
	}
}
