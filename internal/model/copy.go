package model

import "time"

type Copy struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	BookID      int       `json:"book_id" gorm:"not null;index"`
	InventoryNo string    `json:"inventory_no" gorm:"not null;uniqueIndex"`
	Status      string    `json:"status" gorm:"not null;default:available"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
