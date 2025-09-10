package main

import (
	"context"
	"fmt"

	"github.com/feiyizhou/base/clients"
	msgcenter "github.com/feiyizhou/base/msgCenter"
)

func main() {
	msgcenter.NewConsumer(
		context.Background(),
		"topic",
		"group",
		"consumer1",
		50,
		0,
		clients.NewRedisClient(
			context.Background(),
			clients.RedisConf{
				Addr:         "172.16.198.120:32763",
				Password:     "LinCloud1!2@3#4$",
				DB:           0,
				DialTimeout:  5,
				WriteTimeout: 5,
				ReadTimeout:  5,
				PoolTimeout:  5,
			},
		), func(msg msgcenter.Msg) error {
			fmt.Printf("received message, msgType: %s, source: %s, timestamp: %s, msgid: %v\n", msg.MsgType, msg.Source, msg.Timestamp, msg.MsgID)
			return nil
		}).StartConsuming()
}
