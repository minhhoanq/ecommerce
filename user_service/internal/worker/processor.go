package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/user_service/internal/email"
	"github.com/minhhoanq/ecommerce/user_service/internal/usecase/rest/repo"
	"go.uber.org/zap"
)

const (
	QueueCritial = "critial"
	QueueDefault = "default"
)

type TaskProcessor interface {
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	Start() error
	Shutdown()
}

type RedisTaskProcessor struct {
	server *asynq.Server
	l      logger.Interface
	mailer email.EmailSender
	q      repo.Querier
}

func NewRedisTaskProcessor(redisOpts asynq.RedisClientOpt,
	mailer email.EmailSender,
	q repo.Querier,
	l logger.Interface) TaskProcessor {
	server := asynq.NewServer(redisOpts,
		asynq.Config{
			Queues: map[string]int{
				QueueCritial: 10,
				QueueDefault: 5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				l.Error("process task failed", zap.String("Error", err.Error()), zap.String("Type", task.Type()), zap.ByteString("payload", task.Payload()))
			}),
		},
	)

	return &RedisTaskProcessor{
		server: server,
		mailer: mailer,
		q:      q,
		l:      l,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
}

func (processor *RedisTaskProcessor) Shutdown() {
	processor.server.Shutdown()
}
