package repository

import (
	"context"

	"libratrack/internal/model"

	"gorm.io/gorm"
)

type BookRepository interface {
	Create(book model.Book) model.Book
	GetByID(id int) (model.Book, bool)
	Update(id int, newBook model.Book) (model.Book, bool)
	Delete(id int) bool
	List() []model.Book
}

type gormBookRepository struct {
	db *gorm.DB
}

func NewGormBookRepository(db *gorm.DB) BookRepository {
	return &gormBookRepository{db: db}
}

func (r *gormBookRepository) Create(book model.Book) model.Book {
	if r.db == nil {
		return book
	}
	_ = r.db.WithContext(context.Background()).Create(&book).Error
	return book
}

func (r *gormBookRepository) GetByID(id int) (model.Book, bool) {
	var book model.Book
	if r.db == nil {
		return model.Book{}, false
	}
	if err := r.db.WithContext(context.Background()).Preload("Authors").Preload("Category").First(&book, id).Error; err != nil {
		return model.Book{}, false
	}
	return book, true
}

func (r *gormBookRepository) Update(id int, newBook model.Book) (model.Book, bool) {
	if r.db == nil {
		return model.Book{}, false
	}
	newBook.ID = id
	if err := r.db.WithContext(context.Background()).Model(&model.Book{}).Where("id = ?", id).Updates(&newBook).Error; err != nil {
		return model.Book{}, false
	}
	return newBook, true
}

func (r *gormBookRepository) Delete(id int) bool {
	if r.db == nil {
		return false
	}
	return r.db.WithContext(context.Background()).Delete(&model.Book{}, id).Error == nil
}

func (r *gormBookRepository) List() []model.Book {
	var books []model.Book
	if r.db == nil {
		return nil
	}
	if err := r.db.WithContext(context.Background()).Preload("Authors").Preload("Category").Order("id").Find(&books).Error; err != nil {
		return nil
	}
	return books
}
