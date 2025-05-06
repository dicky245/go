package model

import (
	"time"
)

// Jadwal represents a scheduled seminar or presentation
type Jadwal struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	KelompokID uint      `gorm:"column:kelompok_id" json:"kelompok_id"`
	Ruangan    string    `gorm:"column:ruangan" json:"ruangan"`
	Waktu      time.Time `gorm:"column:waktu" json:"waktu"`
	UserID     uint      `gorm:"column:user_id" json:"user_id"` // User who created the jadwal (usually a dosen)
	Penguji1   uint      `gorm:"column:penguji1" json:"penguji1"`
	Penguji2   uint      `gorm:"column:penguji2" json:"penguji2"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
	
	// Virtual fields (not stored in database)
	KelompokNama string `gorm:"-" json:"kelompok_nama,omitempty"`
	Penguji1Nama string `gorm:"-" json:"penguji1_nama,omitempty"`
	Penguji2Nama string `gorm:"-" json:"penguji2_nama,omitempty"`
}

// TableName specifies the table name for Jadwal
func (Jadwal) TableName() string {
	return "jadwal"
}
