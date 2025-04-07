package model

type Role struct {
    ID   uint   `gorm:"primaryKey" json:"id"`
    RoleName string `gorm:"unique;not null" json:"role_name"`
}
func (Role) TableName() string {
    return "roles"
}
