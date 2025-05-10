package model

import (
	"time"
)

// Jadwal represents a scheduled seminar or presentation
type Jadwal struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	KelompokID   uint      `gorm:"column:kelompok_id" json:"kelompok_id"`
	WaktuMulai   time.Time `gorm:"column:waktu_mulai" json:"waktu_mulai"`
	WaktuSelesai time.Time `gorm:"column:waktu_selesai" json:"waktu_selesai"`
	UserID       uint      `gorm:"column:user_id" json:"user_id"` // User who created the jadwal (usually a dosen)
	RuanganID    uint      `gorm:"column:ruangan_id" json:"ruangan_id"`
	KPAID        uint      `gorm:"column:KPA_id" json:"KPA_id"`
	ProdiID      uint      `gorm:"column:prodi_id" json:"prodi_id"`
	TMID         uint      `gorm:"column:TM_id" json:"TM_id"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
	
	// Virtual fields (not stored in database)
	KelompokNama string `gorm:"-" json:"kelompok_nama,omitempty"`
	Ruangan      string `gorm:"-" json:"ruangan,omitempty"`
}

// TableName specifies the table name for Jadwal
func (Jadwal) TableName() string {
	return "jadwal"
}