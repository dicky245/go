package model

import "time"

type Bimbingan struct {
	BimbinganID      uint      `json:"bimbingan_id" gorm:"column:bimbingan_id;primaryKey"`
	Keperluan        string    `json:"keperluan" gorm:"column:keperluan"`
	Deskripsi        string    `json:"deskripsi" gorm:"column:deskripsi"`
	RencanaBimbingan time.Time `json:"rencana_bimbingan" gorm:"column:rencana_bimbingan"`
	Status           string    `json:"status" gorm:"column:status;type:enum('Pending','Diterima','Ditolak','Selesai');default:'Pending'"`
	// DosenID          *uint     `json:"dosen_id" gorm:"column:dosen_id;default:null"`
	// KelompokID       *uint     `json:"kelompok_id" gorm:"column:kelompok_id;default:null"`
	UserID    uint      `json:"user_id" gorm:"column:user_id"` // Pemilik request
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// Override nama tabel agar sesuai dengan database
func (Bimbingan) TableName() string {
	return "request_bimbingan"
}
