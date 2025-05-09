package model

import "time"

type KelompokMahasiswa struct {
    ID         uint      `json:"id" gorm:"column:id;primaryKey"`
    UserID     uint      `json:"user_id" gorm:"column:user_id"`
    KelompokID uint      `json:"kelompok_id" gorm:"column:kelompok_id"`
    CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
    UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
    
    // Relasi
    Kelompok   Kelompok  `gorm:"foreignKey:KelompokID" json:"kelompok"`
}

func (KelompokMahasiswa) TableName() string {
    return "kelompok_mahasiswa"
}