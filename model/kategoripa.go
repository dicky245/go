package model

import "time"

type KategoriPA struct {
    ID         uint      `json:"id" gorm:"column:id;primaryKey"`
    KategoriPA string    `json:"kategori_pa" gorm:"column:kategori_pa"`
    CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
    UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (KategoriPA) TableName() string {
    return "kategori_pa"
}