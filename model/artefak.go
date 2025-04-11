package model

import "time"

// type Artefak struct {
// 	ArtefakID     int       `json:"artefak_id" gorm:"column:artefak_id;primaryKey"`
// 	Judul         string    `json:"judul" gorm:"column:judul"`
// 	Deskripsi     string    `json:"deskripsi" gorm:"column:deskripsi"`
// 	File          string    `json:"file" gorm:"column:file"`
// 	Batas         string    `json:"batas" gorm:"column:batas"`
// 	TanggalSubmit string    `json:"tanggalsubmit" gorm:"column:tanggalsubmit"`
// 	Status        string    `json:"status" gorm:"column:status;type:enum('Pending','Submitted','Approved','Rejected');default:'Pending'"`
// 	UserID        uint      `json:"user_id" gorm:"column:user_id"`
// 	User          User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
// 	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
// 	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
// }

type Artefak struct {
	ArtefakID uint      `json:"artefak_id" gorm:"column:artefak_id;primaryKey"`
	SubmitID  uint      `gorm:"column:submit_id" json:"submit_id"`
	File      string    `json:"file" gorm:"column:file"`
	UserID    uint      `json:"user_id" gorm:"column:user_id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Submit    Submit    `gorm:"foreignKey:SubmitID" json:"pengumpulan,omitempty"`
}

func (Artefak) TableName() string {
	return "artefaks"
}
