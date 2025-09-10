package msgcenter

import (
	"context"

	"github.com/feiyizhou/base/clients"
	"github.com/feiyizhou/base/logger"
)

type Producer struct {
	topic string
	group string
	ctx   context.Context
	redis *clients.RedisClient
}

func NewProducer(ctx context.Context, redis *clients.RedisClient, topic, group string) *Producer {
	return &Producer{
		ctx:   ctx,
		topic: topic,
		group: group,
		redis: redis,
	}
}

func (p *Producer) Publish(msg map[string]any) (string, error) {
	var (
		err   error
		msgId string
	)
	if msgId, err = p.redis.PublishToStream(p.ctx, p.topic, msg); err != nil {
		logger.Errorf("publish msg to topic %s failed, err: %v", p.topic, err)
		return "", err
	}
	logger.Infof("publish msg to topic %s success, msgId: %s", p.topic, msgId)
	return msgId, nil
}
