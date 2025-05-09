package model

import (
	"time"
)

// Pengumuman represents an announcement in the system
type Pengumuman struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Judul            string    `json:"judul" gorm:"column:judul"`
	Deskripsi        string    `json:"deskripsi" gorm:"column:deskripsi"`
	TanggalPenulisan time.Time `json:"tanggal_penulisan" gorm:"column:tanggal_penulisan"`
	File             string    `json:"file" gorm:"column:file"`
	Status           string    `json:"status" gorm:"column:status;type:enum('aktif','non-aktif')"`
	UserID           uint      `json:"user_id" gorm:"column:user_id"`
	KPAID            uint      `json:"kpa_id" gorm:"column:KPA_id"`
	ProdiID          uint      `json:"prodi_id" gorm:"column:prodi_id"`
	TMID             uint      `json:"tm_id" gorm:"column:TM_id"`
	CreatedAt        time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// TableName specifies the table name for the Pengumuman model
func (Pengumuman) TableName() string {
	return "pengumuman"
}