package consumers

import (
	"context"

	"github.com/google/uuid"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/notification_service/internal/dataaccess/database"
	"github.com/minhhoanq/ecommerce/notification_service/internal/email"
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
	emailSender              email.EmailSender
	notificationDataAccessor database.NotificationDataAccessor
	l                        logger.Interface
}

func NewEmailNotifyCompletedMessageHandler(
	emailSender email.EmailSender,
	notificationDataAccessor database.NotificationDataAccessor,
	l logger.Interface,
) EmailNotifyCompletedMessageHandler {
	return &emailNotifyCompletedMessageHandler{
		emailSender:              emailSender,
		notificationDataAccessor: notificationDataAccessor,
		l:                        l,
	}
}

func (e emailNotifyCompletedMessageHandler) Handle(ctx context.Context, message EmailNotifyCompleted) error {
	// e.emailSender.SendMail()
	return nil
}
