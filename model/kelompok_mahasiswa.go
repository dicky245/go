package model

import (
    "time"
)

type KelompokMahasiswa struct {
    ID         uint     `gorm:"primaryKey" json:"id"`
    UserID     uint     `json:"user_id" gorm:"column:user_id"`
    KelompokID uint     `json:"kelompok_id" gorm:"column:kelompok_id"`
    Kelompok   Kelompok `gorm:"foreignKey:KelompokID" json:"kelompok"` // Relasi ke tabel Kelompok
    CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
    UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (KelompokMahasiswa) TableName() string {
    return "kelompok_mahasiswa"
}