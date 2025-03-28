package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationType
type NotificationType string

const (
	SystemNotification NotificationType = "SYSTEM"
	EmailNotification  NotificationType = "EMAIL"
)

// NotificationStatus
type NotificationStatus string

const (
	Unread NotificationStatus = "UNREAD"
	Read   NotificationStatus = "READ"
	Sent   NotificationStatus = "SENT"
	FAILED NotificationStatus = "FAILED"
)

type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Type      string             `bson:"type" json:"type"`
	Title     string             `bson:"title" json:"title"`
	Message   string             `bson:"message" json:"message"`
	Status    string             `bson:"status" json:"status"`
	Metadata  map[string]string  `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type SendNotificationRequest struct {
	UserID   string            `bson:"user_id" json:"user_id"`
	Type     string            `bson:"type" json:"type"`
	Title    string            `bson:"title" json:"title"`
	Message  string            `bson:"message" json:"message"`
	Status   string            `bson:"status" json:"status"`
	Metadata map[string]string `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

type SendNotificationResponse struct {
	ID     string `bson:"_id,omitempty" json:"id"`
	Status string `bson:"status" json:"status"`
}
