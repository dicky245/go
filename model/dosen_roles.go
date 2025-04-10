package model

type DosenRole struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	UserID     uint   `json:"user_id"`
	NamaDosen  string `json:"nama_dosen"`
	RoleID     uint   `json:"role_id"`
	NamaRole   string `json:"nama_role"`
	Prodi      string `json:"prodi"`
	JenisPA    string   `json:"jenis_pa"`
}

func (DosenRole) TableName() string {
    return "dosen_roles"
}
