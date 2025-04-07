package model

type DosenRole struct {
    ID     uint `gorm:"primaryKey" json:"id"`
    UserID uint `gorm:"not null" json:"user_id"`
    Nama string `gorm:"not null" json:"nama_dosen"`
    Prodi   string `gorm:"not null" json:"prodi"` 
    Tingkat uint   `gorm:"not null" json:"tingkat"`
    RoleID uint `gorm:"not null" json:"role_id"` 
    NamaRole string `gorm:"not null" json:"nama_role"`
}

func (DosenRole) TableName() string {
    return "dosen_roles"
}
