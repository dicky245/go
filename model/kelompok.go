package model

type Kelompok struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Nomor     string `gorm:"type:varchar(100);not null";index:idx_unique_kelompok,unique" json:"nomor"`
	JenisPA  string `gorm:"type:varchar(10);not null";index:idx_unique_kelompok,unique" json:"jenis_pa"` 
	Prodi   string `gorm:"type: varchar(100);not null";index:idx_unique_kelompok,unique" json:"Prodi"`
	Angkatan int    `gorm:"not null";index:idx_unique_kelompok,unique" json:"angkatan"`
}
func (Kelompok) TableName() string{
	return "kelompok"
}