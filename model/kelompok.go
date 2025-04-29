package model

import "time"

type Kelompok struct {
    ID             uint      `json:"id" gorm:"column:id;primaryKey"`
    NomorKelompok string    `json:"nomor_kelompok" gorm:"column:nomor_kelompok"`
    KPAID          uint      `json:"kpa_id" gorm:"column:KPA_id"` // Note the capital KPA
    ProdiID        uint      `json:"prodi_id" gorm:"column:prodi_id"`
    TAID           uint      `json:"ta_id" gorm:"column:TA_id"` // Note the capital TA
    Status         string    `json:"status" gorm:"column:status"`
    CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`
    UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at"`

    // Relasi
    Prodi      Prodi      `gorm:"foreignKey:ProdiID" json:"prodi"`
    TahunAjaran TahunAjaran `gorm:"foreignKey:TAID" json:"tahun_ajaran"`
}

func (Kelompok) TableName() string {
    return "kelompok"
}