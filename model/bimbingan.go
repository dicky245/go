package model

import "time"

type Bimbingan struct {
    ID             uint      `json:"id" gorm:"column:id;primaryKey"`
    KelompokID     uint      `json:"kelompok_id" gorm:"column:kelompok_id"`
    UserID         uint      `json:"user_id" gorm:"column:user_id"`
    Keperluan      string    `json:"keperluan" gorm:"column:keperluan"`
    RencanaMulai   time.Time `json:"rencana_mulai" gorm:"column:rencana_mulai"`
    RencanaSelesai time.Time `json:"rencana_selesai" gorm:"column:rencana_selesai"`
    RuanganID      uint      `json:"ruangan_id" gorm:"column:ruangan_id"`
    Status         string    `json:"status" gorm:"column:status;type:enum('menunggu','selesai','disetujui','ditolak');default:'menunggu'"`
    HasilBimbingan string    `json:"hasil_bimbingan" gorm:"column:hasil_bimbingan"`

    CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`

    // Relasi opsional
    Kelompok KelompokMahasiswa `gorm:"foreignKey:KelompokID" json:"kelompok,omitempty"`
    Ruangan  Ruangan           `gorm:"foreignKey:RuanganID" json:"ruangan,omitempty"`
}


func (Bimbingan) TableName() string {
    return "request_bimbingan"
}
