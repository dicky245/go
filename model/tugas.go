package model

import "time"

type Tugas struct {
    ID                uint            `json:"id" gorm:"column:id;primaryKey"`
    UserID            uint            `json:"user_id" gorm:"column:user_id"`
    JudulTugas        string          `json:"judul_tugas" gorm:"column:Judul_Tugas"`
    DeskripsiTugas    string          `json:"deskripsi_tugas" gorm:"column:Deskripsi_Tugas"`
    KPAID             uint            `json:"kpa_id" gorm:"column:KPA_id"`
    ProdiID           uint            `json:"prodi_id" gorm:"column:prodi_id"`
    TMID              uint            `json:"tm_id" gorm:"column:TM_id"`
    TanggalPengumpulan time.Time      `json:"tanggal_pengumpulan" gorm:"column:tanggal_pengumpulan"`
    File              string          `json:"file" gorm:"column:file"`
    Status            string          `json:"status" gorm:"column:status;default:'berlangsung'"`
    KategoriTugas     string          `json:"kategori_tugas" gorm:"column:kategori_tugas;default:'Tugas'"`
    CreatedAt         time.Time       `json:"created_at" gorm:"column:created_at;autoCreateTime"`
    UpdatedAt         time.Time       `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`

    // Relasi
    Prodi             Prodi           `gorm:"foreignKey:ProdiID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"prodi,omitempty"`
    KategoriPA        KategoriPA      `gorm:"foreignKey:KPAID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"kategori_pa,omitempty"`
    TahunMasuk        TahunAjaran      `gorm:"foreignKey:TMID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"tahun_masuk,omitempty"`
    
    PengumpulanTugas  []PengumpulanTugas `gorm:"foreignKey:TugasID" json:"pengumpulan_tugas,omitempty"`
}

func (Tugas) TableName() string {
    return "tugas"
}