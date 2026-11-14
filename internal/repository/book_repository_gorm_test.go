package repository

import (
	"testing"

	"libratrack/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGormBookRepositoryPersistsAuthorBio(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.Book{}, &model.Author{}, &model.Category{}); err != nil {
		t.Fatalf("migrate schema: %v", err)
	}

	repo := NewGormBookRepository(db)
	originalBio := "Original bio."
	book := model.Book{
		Title:    "The Hound of the Baskervilles",
		ISBN:     "978-1-2345-6789-0",
		Category: model.Category{Name: "Classics"},
		Authors:  []model.Author{{FullName: "Arthur Conan Doyle", Bio: &originalBio}},
	}

	created := repo.Create(book)
	if len(created.Authors) != 1 || created.Authors[0].Bio == nil || *created.Authors[0].Bio != originalBio {
		t.Fatalf("expected created author bio to be persisted, got %#v", created.Authors)
	}

	updatedBio := "Updated bio."
	created.Authors[0].Bio = &updatedBio
	if _, ok := repo.Update(created.ID, created); !ok {
		t.Fatal("expected update to succeed")
	}

	stored, ok := repo.GetByID(created.ID)
	if !ok {
		t.Fatal("expected book to exist after update")
	}
	if len(stored.Authors) != 1 || stored.Authors[0].Bio == nil || *stored.Authors[0].Bio != updatedBio {
		t.Fatalf("expected updated author bio in stored book, got %#v", stored.Authors)
	}
}
