package model

import "time"

// Penguji represents an examiner assigned to a kelompok
type Penguji struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"column:user_id" json:"user_id"`
	KelompokID uint      `gorm:"column:kelompok_id" json:"kelompok_id"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName specifies the table name for Penguji
func (Penguji) TableName() string {
	return "penguji"
}