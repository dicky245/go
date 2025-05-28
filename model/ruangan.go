package model

import "time"

// Ruangan represents a room where seminars can be held
type Ruangan struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Ruangan   string    `gorm:"column:ruangan" json:"ruangan"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName specifies the table name for Ruangan
func (Ruangan) TableName() string {
	return "ruangan"
}