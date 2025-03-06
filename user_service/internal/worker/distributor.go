package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/minhhoanq/ecommerce/common/logger"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedistaskDistributor struct {
	client *asynq.Client
	l      logger.Interface
}

func NewRedisTaskDistributor(l logger.Interface, redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedistaskDistributor{
		client: client,
		l:      l,
	}
}
