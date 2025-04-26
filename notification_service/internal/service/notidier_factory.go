package service

import (
	"fmt"

	"github.com/minhhoanq/ecommerce/notification_service/internal/email"
)

type NotifierFactory interface {
	GetNotifier(t NotificationType) (Notifier, error)
}

type notifierFactory struct {
	emailSender email.EmailSender
}

func NewNotifierFactory(emailSender email.EmailSender) NotifierFactory {
	return &notifierFactory{
		emailSender: emailSender,
	}
}

func (f *notifierFactory) GetNotifier(t NotificationType) (Notifier, error) {
	switch t {
	case TypeEmail:
		return &EmailNotifier{
			emailSender: f.emailSender,
		}, nil
	case TypeSMS:
		return &SMSNotifier{}, nil
	default:
		return nil, fmt.Errorf("no notifier type: %s", t)
	}
}
