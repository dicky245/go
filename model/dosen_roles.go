package model

type DosenRole struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	UserID     uint   `json:"user_id"`
	NamaDosen  string `json:"nama_dosen"`
	RoleID     uint   `json:"role_id"`
	NamaRole   string `json:"nama_role"`
	Prodi      string `json:"prodi"`
	Tingkat    uint   `json:"tingkat"`
}

func (DosenRole) TableName() string {
    return "dosen_roles"
}
