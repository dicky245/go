package model

import "time"

type Jadwal struct {
	JadwalID  uint      `json:"jadwal_id" gorm:"column:jadwal_id;primaryKey"`
	Tanggal   string    `json:"tanggal" gorm:"column:tanggal"`
	Ruangan   string    `json:"ruangan" gorm:"column:ruangan"`
	Jam       string    `json:"jam" gorm:"column:jam"`
	Kelompok  int       `json:"kelompok" gorm:"column:kelompok"`
	UserID    uint      `json:"user_id" gorm:"column:user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Jadwal) TableName() string {
	return "jadwal"
}
