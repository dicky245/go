package model

import (
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	TypeJadwal      NotificationType = "jadwal"
	TypePengumuman  NotificationType = "pengumuman"
	TypeSubmission  NotificationType = "submission"
)
		
// Notification represents a notification in the system
type Notification struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"column:title"`
	Message     string         `json:"message" gorm:"column:message"`
	Type        NotificationType `json:"type" gorm:"column:type;type:enum('jadwal','pengumuman','submission')"`
	ReferenceID uint           `json:"reference_id" gorm:"column:reference_id"` // ID of the related item (jadwal, pengumuman, etc.)
	UserID      uint           `json:"user_id" gorm:"column:user_id"`           // Recipient user ID
	IsRead      bool           `json:"is_read" gorm:"column:is_read;default:false"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"column:updated_at"`
}

// TableName specifies the table name for the Notification model
func (Notification) TableName() string {
	return "notifications"
}

// DeviceToken represents a user's device token for push notifications
type DeviceToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"column:user_id;uniqueIndex:idx_user_token"`
	Token     string    `json:"token" gorm:"column:token;uniqueIndex:idx_user_token"`
	Platform  string    `json:"platform" gorm:"column:platform"` // "android" or "ios"
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// TableName specifies the table name for the DeviceToken model
func (DeviceToken) TableName() string {
	return "device_tokens"
}
