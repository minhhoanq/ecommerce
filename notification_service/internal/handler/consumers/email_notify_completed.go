package consumers

import (
	"context"

	"github.com/google/uuid"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/notification_service/internal/service"
)

const (
	TOPIC_NAME_NOTIFICATION_SERVICE_EMAIL_NOTIFY_COMPLETED = "notification_service_email_notify_completed"
)

type EmailNotifyCompleted struct {
	UserID uuid.UUID
	Type   string
}

type EmailNotifyCompletedMessageHandler interface {
	Handle(ctx context.Context, message EmailNotifyCompleted) error
}

type emailNotifyCompletedMessageHandler struct {
	notificationService service.NotificationService
	l                   logger.Interface
}

func NewEmailNotifyCompletedMessageHandler(
	notificationService service.NotificationService,
	l logger.Interface,
) EmailNotifyCompletedMessageHandler {
	return &emailNotifyCompletedMessageHandler{
		notificationService: notificationService,
		l:                   l,
	}
}

func (e emailNotifyCompletedMessageHandler) Handle(ctx context.Context, message EmailNotifyCompleted) error {
	return e.notificationService.SendEmailWhenSignup(ctx)
}
