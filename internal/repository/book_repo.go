// Package repository holds the storage interface and its in-memory implementation.
package repository

import (
	"sync"
	"homework/internal/model"
	"github.com/google/uuid"
)

type BookRepository struct {
	mutex    sync.RWMutex
	books map[string]model.Book
}

func NewBookRepository() *BookRepository {
	return &BookRepository{
		books: make(map[string]model.Book),
	}
}

// приклад створення книги для репозиторію
func (bookRepo *BookRepository) Create(book model.Book) model.Book {
	bookRepo.mutex.Lock()
	defer bookRepo.mutex.Unlock()
	
	book.ID = uuid.New().String()
	bookRepo.books[book.ID] = book
	return book
}

// протрібно ще зробити методи GetAll, GetByID, Update, Delete для CRUD

