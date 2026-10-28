package model

import "time"

type Book struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title         string    `json:"title" gorm:"not null"`
	ISBN          string    `json:"isbn" gorm:"not null;uniqueIndex"`
	CategoryID    int       `json:"category_id" gorm:"not null;index"`
	Category      Category  `json:"category" gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT"`
	Authors       []Author  `json:"authors" gorm:"many2many:book_authors;constraint:OnDelete:CASCADE"`
	PublishedYear *int      `json:"published_year"`
	Description   *string   `json:"description"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
