package msgcenter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/feiyizhou/base/clients"

	"github.com/redis/go-redis/v9"
)

type Consumer struct {
	topic       string
	group       string
	count       int64
	block       int64
	consumerID  string
	ctx         context.Context
	redis       *clients.RedisClient
	handlerFunc func(msg Msg) error
}

func NewConsumer(ctx context.Context, topic, group, consumerID string, count, block int64, redis *clients.RedisClient, handlerFunc func(msg Msg) error) *Consumer {
	if count <= 0 {
		count = 1
	}
	if block <= 0 {
		block = 0
	}
	return &Consumer{
		ctx:         ctx,
		topic:       topic,
		group:       group,
		redis:       redis,
		consumerID:  consumerID,
		handlerFunc: handlerFunc,
	}
}

func (c *Consumer) StartConsuming() {
	var (
		err   error
		xmsgs []redis.XMessage
	)
	if err = c.redis.CreateConsumerGroup(c.ctx, c.topic, c.group); err != nil {
		fmt.Printf("create consumer group %s for topic %s failed, err: %v\n", c.group, c.topic, err)
		panic(err)
	}
	for {
		if xmsgs, err = c.redis.ReadFromStream(c.ctx, c.topic, c.group, c.consumerID, c.count, c.block); err != nil {
			continue
		}
		for _, xmsg := range xmsgs {
			var msg Msg
			xmsgBytes, _ := json.Marshal(xmsg.Values)
			fmt.Printf("xmsgBytes: %s\n", string(xmsgBytes))
			if err = json.Unmarshal(xmsgBytes, &msg); err != nil {
				continue
			}
			if err = c.handlerFunc(msg); err != nil {
				continue
			}
			if err = c.redis.AckMessage(c.ctx, c.topic, c.group, xmsg.ID); err != nil {
				continue
			}
		}
	}
}
