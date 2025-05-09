package model

import "time"

type Kelompok struct {
    ID          uint      `json:"id" gorm:"column:id;primaryKey"`
    NomorKelompok string  `json:"nomor_kelompok" gorm:"column:nomor_kelompok"`
    KPAID       uint      `json:"kpa_id" gorm:"column:KPA_id"`
    ProdiID     uint      `json:"prodi_id" gorm:"column:prodi_id"`
    TMID        uint      `json:"tm_id" gorm:"column:TM_id"`
    Status      string    `json:"status" gorm:"column:status"`
    CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Kelompok) TableName() string {
    return "kelompok"
}