package model

import "time"

type Prodi struct {
    ID          uint      `json:"id" gorm:"column:id;primaryKey"`
    NamaProdi   string    `json:"nama_prodi" gorm:"column:nama_prodi"`
    MaksProject int       `json:"maks_project" gorm:"column:maks_project"`
    CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Prodi) TableName() string {
    return "prodi"
}