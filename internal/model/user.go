package model

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Email     string    `json:"email" gorm:"not null;uniqueIndex"`
	Role      string    `json:"role" gorm:"not null;default:user"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
