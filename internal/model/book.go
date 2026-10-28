package model

import (
	"encoding/json"
	"strings"
	"time"
)

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

func (b Book) LegacyAuthor() string {
	if len(b.Authors) > 0 {
		if name := strings.TrimSpace(b.Authors[0].FullName); name != "" {
			return name
		}
	}
	return ""
}

func (b Book) LegacyCategory() string {
	if name := strings.TrimSpace(b.Category.Name); name != "" {
		return name
	}
	return ""
}

func (b Book) MarshalJSON() ([]byte, error) {
	type bookAlias struct {
		ID            int       `json:"id"`
		Title         string    `json:"title"`
		ISBN          string    `json:"isbn"`
		Author        string    `json:"author,omitempty"`
		Category      string    `json:"category,omitempty"`
		PublishedYear *int      `json:"published_year"`
		Description   *string   `json:"description"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}

	aux := bookAlias{
		ID:            b.ID,
		Title:         b.Title,
		ISBN:          b.ISBN,
		PublishedYear: b.PublishedYear,
		Description:   b.Description,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
	}
	if auth := b.LegacyAuthor(); auth != "" {
		aux.Author = auth
	}
	if cat := b.LegacyCategory(); cat != "" {
		aux.Category = cat
	}
	return json.Marshal(aux)
}

func (b *Book) UnmarshalJSON(data []byte) error {
	type rawBook struct {
		ID            int       `json:"id"`
		Title         string    `json:"title"`
		ISBN          string    `json:"isbn"`
		Author        string    `json:"author"`
		Category      string    `json:"category"`
		CategoryID    int       `json:"category_id"`
		AuthorIDs     []int     `json:"author_ids"`
		PublishedYear *int      `json:"published_year"`
		Description   *string   `json:"description"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}

	var aux rawBook
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	b.ID = aux.ID
	b.Title = aux.Title
	b.ISBN = aux.ISBN
	b.PublishedYear = aux.PublishedYear
	b.Description = aux.Description
	b.CreatedAt = aux.CreatedAt
	b.UpdatedAt = aux.UpdatedAt

	if aux.CategoryID > 0 {
		b.CategoryID = aux.CategoryID
		b.Category = Category{ID: aux.CategoryID}
	}
	if strings.TrimSpace(aux.Category) != "" {
		b.Category = Category{Name: aux.Category}
		if b.CategoryID == 0 {
			b.CategoryID = 0
		}
	}

	if len(aux.AuthorIDs) > 0 {
		b.Authors = make([]Author, 0, len(aux.AuthorIDs))
		for _, id := range aux.AuthorIDs {
			b.Authors = append(b.Authors, Author{ID: id})
		}
	} else if strings.TrimSpace(aux.Author) != "" {
		b.Authors = []Author{{FullName: aux.Author}}
	}

	return nil
}
