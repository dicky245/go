package model

import "time"

type Device_Token struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      int    `json:"user_id" gorm:"column:user_id;unique"`
	TokenDevice string `json:"token_device" gorm:"column:token_device"`
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

func (Device_Token) TableName() string {
	return "device_token"
}
