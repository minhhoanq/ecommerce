package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/notification_service/internal/dataaccess/database"
	"github.com/minhhoanq/ecommerce/notification_service/internal/email"
	pb "github.com/minhhoanq/ecommerce/notification_service/internal/generated/notification_service"
	"github.com/minhhoanq/ecommerce/notification_service/internal/generated/user_service"
)

const path = "http://localhost:8000/v1/verify_email"

type NotificationService interface {
	SendNotification(ctx context.Context, arg *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error)
	SendEmailWhenSignup(ctx context.Context, arg *SendEmailWhenSignupRequest, notiType NotificationType) error
}

type notificationService struct {
	l                        logger.Interface
	notificationDataAccessor database.NotificationDataAccessor
	userServiceGRPCClient    user_service.UserServiceClient
	notifier                 Notifier
	emailSender              email.EmailSender
}

func NewNotificationService(
	l logger.Interface,
	notificationDataAccessor database.NotificationDataAccessor,
	userServiceGRPCClient user_service.UserServiceClient,
	emailSender email.EmailSender,
) NotificationService {

	return &notificationService{
		l:                        l,
		notificationDataAccessor: notificationDataAccessor,
		userServiceGRPCClient:    userServiceGRPCClient,
		emailSender:              emailSender,
	}
}

func (n *notificationService) setNotifier(notifier Notifier) {
	n.notifier = notifier
}

func (n *notificationService) SendNotification(ctx context.Context, arg *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	payload := &database.SendNotificationRequest{
		UserID:   arg.UserId,
		Type:     arg.Type.String(),
		Title:    arg.Title,
		Message:  arg.Message,
		Status:   arg.Status.String(),
		Metadata: arg.Metadata,
	}

	notification, err := n.notificationDataAccessor.CreateNotification(ctx, payload)
	if err != nil {
		return nil, err
	}

	response := &pb.SendNotificationResponse{
		NotificationId: notification.ID,
		Status:         notification.Status,
	}

	return response, nil
}

type SendEmailWhenSignupRequest struct {
	UserID        uuid.UUID
	Email         string
	Username      string
	SecretCode    string
	VerifyEmailID int
}

func (n *notificationService) SendEmailWhenSignup(ctx context.Context, arg *SendEmailWhenSignupRequest, notiType NotificationType) error {
	subject := "Welcome to LIFEAT"
	verifyUrl := fmt.Sprintf("%s?email_id=%d&secret_code=%s", path, arg.VerifyEmailID, arg.SecretCode)
	content := fmt.Sprintf(`Hello %s, <br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">Click here<a> to verify your email address.<br/>`, arg.Username, verifyUrl)
	to := []string{arg.Email}

	payload := &database.SendNotificationRequest{
		UserID:  arg.UserID.String(),
		Type:    "Email",
		Title:   subject,
		Message: content,
		Status:  "Success",
		Metadata: map[string]string{
			"url": verifyUrl,
		},
	}
	_, err := n.notificationDataAccessor.CreateNotification(ctx, payload)
	if err != nil {
		return err
	}

	n.setNotifier(&EmailNotifier{emailSender: n.emailSender})
	if err != nil {
		return err
	}

	notiParams := &NotifierParams{
		notiType:    notiType,
		subject:     subject,
		content:     content,
		to:          to,
		cc:          nil,
		bcc:         nil,
		attachFiles: nil,
	}

	return n.notifier.Send(ctx, notiParams)
}
