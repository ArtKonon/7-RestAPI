// Package repository holds the storage interface and its in-memory implementation.
package repository

import (
    "sort"
    "sync"
    "time"

    "homework/internal/model"
)

type BookRepository struct {
	mutex     sync.RWMutex
	books  map[int]model.Book
	nextID int
}

func NewBookRepository() *BookRepository {
	return &BookRepository{
		books:  make(map[int]model.Book),
		nextID: 1,
	}
}

func (repo *BookRepository) Create(book model.Book) model.Book {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	book.ID = repo.nextID
	now := time.Now().UTC()
	book.CreatedAt = now
	book.UpdatedAt = now

	repo.books[book.ID] = book
	repo.nextID++

	return book
}

func (repo *BookRepository) GetByID(id int) (model.Book, bool) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	book, ok := repo.books[id]
	return book, ok
}

func (repo *BookRepository) List() []model.Book {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	ids := make([]int, 0, len(repo.books))
	for id := range repo.books {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	out := make([]model.Book, 0, len(ids))
	for _, id := range ids {
		out = append(out, repo.books[id])
	}
	return out
}

func (r *BookRepository) Update(id int, newBook model.Book) (model.Book, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	old, ok := r.books[id]
	if !ok {
		return model.Book{}, false
	}

	newBook.ID = id
	newBook.CreatedAt = old.CreatedAt
	newBook.UpdatedAt = time.Now().UTC()

	r.books[id] = newBook
	return newBook, true
}

func (r *BookRepository) Delete(id int) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.books[id]; !ok {
		return false
	}

	delete(r.books, id)
	return true
}

