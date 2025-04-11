package model

import "time"

type Pengumuman struct {
	NotifikasiID int       `json:"notifikasi_id" gorm:"column: notifikasi_id;primaryKey;autoIncrement"`
	Judul        string    `json:"judul" gorm:"column:judul"`
	Pesan        string    `json:"pesan" gorm:"column:pesan"`
	Status       string    `json:"status" gorm:"column:status;type:enum('Unread','Read');default:'Unread'"`
	UserID       uint      `json:"user_id" gorm:"column:user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Pengumuman) TableName() string {
	return "pengumumen"
}
