package model

import "time"

type PengumpulanTugas struct {
    ID            uint      `json:"id" gorm:"column:id;primaryKey"`
    KelompokID    uint      `json:"kelompok_id" gorm:"column:kelompok_id"`
    TugasID       uint      `json:"tugas_id" gorm:"column:tugas_id"`
    WaktuSubmit   time.Time `json:"waktu_submit" gorm:"column:waktu_submit"`
    FilePath      string    `json:"file_path" gorm:"column:file_path"`
    Status        string    `json:"status" gorm:"column:status;default:'Belum'"`
    CreatedAt     time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
    UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
    Tugas         Tugas     `json:"tugas" gorm:"foreignKey:TugasID;references:ID"`
}

func (PengumpulanTugas) TableName() string {
    return "pengumpulan_tugas"
}
