package model

import "time"

type Bimbingan struct {
    ID             uint      `json:"id" gorm:"column:id;primaryKey"`
    KelompokID     uint      `json:"kelompok_id" gorm:"column:kelompok_id"` // ID kelompok
    UserID         uint      `json:"user_id" gorm:"column:user_id"`         // ID user dari API
    Keperluan      string    `json:"keperluan" gorm:"column:keperluan"`
    RencanaMulai   time.Time `json:"rencana_mulai" gorm:"column:rencana_mulai"`
    RencanaSelesai time.Time `json:"rencana_selesai" gorm:"column:rencana_selesai"`
    Lokasi         string    `json:"lokasi" gorm:"column:lokasi"`
    Status         string    `json:"status" gorm:"column:status;type:enum('menunggu','selesai','disetujui','ditolak');default:'menunggu'"`

    CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`

    // Hapus relasi ke User karena ambil dari API
    // User User `gorm:"foreignKey:UserID" json:"user,omitempty"`

    // Relasi ke Kelompok masih boleh jika tabelnya ada
    Kelompok KelompokMahasiswa `gorm:"foreignKey:KelompokID" json:"kelompok,omitempty"`
}

func (Bimbingan) TableName() string {
    return "request_bimbingan"
}
