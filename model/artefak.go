package model

import "time"

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
