package model

import "time"

type Author struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	FullName  string    `json:"full_name" gorm:"not null;index"`
	Bio       *string   `json:"bio"`
	Books     []Book    `json:"books" gorm:"many2many:book_authors;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
