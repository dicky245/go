package model

import "time"

type Submit struct {
	SubmitID  uint      `json:"submit_id" gorm:"column:submit_id;primaryKey"`
	Judul     string    `json:"judul" gorm:"column:judul"`
	Instruksi string    `json:"instruksi" gorm:"column:instruksi"`
	File      string    `json:"file" gorm:"column:file"`
	Batas     string    `json:"batas" gorm:"column:batas"`
	UserID    uint      `json:"user_id" gorm:"column:user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	Artefaks  []Artefak `gorm:"foreignKey:SubmitID" json:"artefaks,omitempty"`
}

func (Submit) TableName() string {
	return "submits"
}
