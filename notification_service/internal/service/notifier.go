package service

import (
	"context"
	"fmt"

	"github.com/minhhoanq/ecommerce/notification_service/internal/email"
)

type NotificationType string

const (
	TypeEmail  NotificationType = "email"
	TypeSMS    NotificationType = "sms"
	TypeSystem NotificationType = "system"
)

type NotifierParams struct {
	notiType NotificationType

	// params email sender
	subject     string
	content     string
	to          []string
	cc          []string
	bcc         []string
	attachFiles []string
}

type Notifier interface {
	Send(ctx context.Context, arg *NotifierParams) error
}

type EmailNotifier struct {
	emailSender email.EmailSender
}

func (e EmailNotifier) Send(ctx context.Context, n *NotifierParams) error {
	return e.emailSender.SendMail(n.subject, n.content, n.to, n.cc, n.bcc, n.attachFiles)
}

type SMSNotifier struct{}

func (e SMSNotifier) Send(ctx context.Context, n *NotifierParams) error {
	fmt.Println("sms notifier send")
	return nil
}
