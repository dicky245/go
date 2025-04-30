package model

import "time"

type TahunAjaran struct {
    ID          uint      `json:"id" gorm:"column:id;primaryKey"`
    TahunAjaran string    `json:"tahun_ajaran" gorm:"column:Tahun_Masuk"` // Note capital T and A
    Status      string    `json:"status" gorm:"column:Status"` // Note capital S
    CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (TahunAjaran) TableName() string { 
    return "tahun_masuk"
}