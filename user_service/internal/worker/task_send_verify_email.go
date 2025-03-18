package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/minhhoanq/ecommerce/user_service/internal/dataaccess/database"
	"github.com/minhhoanq/ecommerce/user_service/internal/util"
	"go.uber.org/zap"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	UserId uuid.UUID `json:"user_id"`
}

func (distributor *RedistaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	distributor.l.Info("Distributed send verify email")
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed marshal payload %w", err)
	}

	distributor.l.Info("Create task send verify email")

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)

	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to queue task %w", err)
	}

	distributor.l.Info("enqueued task", zap.String("Task type", task.Type()),
		zap.ByteString("payload", task.Payload()),
		zap.String("Queue", info.Queue),
		zap.Int("max_retry", info.MaxRetry))

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	processor.l.Info("processing task send verify email", zap.String("type", task.Type()))

	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.userDataAccessor.GetUserByID(ctx, payload.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found %w", asynq.SkipRetry)
		}

		return fmt.Errorf("failed to get user: %w", asynq.SkipRetry)
	}

	verifyEmail, err := processor.userDataAccessor.CreateVerifyEmail(ctx, database.CreateVerifyEmailParams{
		UserId:     user.ID,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})

	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	subject := "Welcome to LIFEAT"
	verifyUrl := fmt.Sprintf("http://localhost:8000/v1/verify_email?email_id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s, <br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">Click here<a> to verify your email address.<br/>`, user.Username, verifyUrl)
	to := []string{user.Email}

	err = processor.mailer.SendMail(subject, content, to, nil, nil, nil)

	processor.l.Info("process task send verify email", zap.String("type", task.Type()), zap.String("email", user.Email))

	return nil
}
