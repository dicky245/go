package model

import "time"

type Tugas struct {
    ID                uint            `json:"id" gorm:"column:id;primaryKey"`
    UserID            uint            `json:"user_id" gorm:"column:user_id"`
    JudulTugas        string          `json:"judul_tugas" gorm:"column:Judul_Tugas"` // Note capital J
    DeskripsiTugas    string          `json:"deskripsi_tugas" gorm:"column:Deskripsi_Tugas"` // Note capital D
    KPAID             uint            `json:"kpa_id" gorm:"column:KPA_id"` // Note capital KPA
    ProdiID           uint            `json:"prodi_id" gorm:"column:prodi_id"`
    TAID              uint            `json:"ta_id" gorm:"column:TA_id"` // Note capital TA
    TanggalPengumpulan time.Time       `json:"tanggal_pengumpulan" gorm:"column:tanggal_pengumpulan"`
    File              string          `json:"file" gorm:"column:file"`
    Status            string          `json:"status" gorm:"column:status;default:'berlangsung'"`
    CreatedAt         time.Time       `json:"created_at" gorm:"column:created_at;autoCreateTime"`
    UpdatedAt         time.Time       `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`

    // Relasi
    Prodi             Prodi           `gorm:"foreignKey:ProdiID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"prodi"`
    KategoriPA        KategoriPA      `gorm:"foreignKey:KPAID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"kategori_pa"`
    TahunAjaran       TahunAjaran     `gorm:"foreignKey:TAID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"tahun_ajaran"`

    PengumpulanTugas  []PengumpulanTugas `gorm:"foreignKey:TugasID" json:"pengumpulan_tugas,omitempty"`
}

func (Tugas) TableName() string {
    return "tugas"
}