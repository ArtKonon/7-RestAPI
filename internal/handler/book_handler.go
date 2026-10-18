// Package handler holds the HTTP handlers for your resource.
package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"homework/internal/model"
	"homework/internal/repository"
)

type BookHandler struct {
	bookRepo *repository.BookRepository
}

func NewBookHandler(repo *repository.BookRepository) *BookHandler {
	return &BookHandler{bookRepo: repo}
}

// CreateBook обробляє POST /books
func (handler *BookHandler) CreateBook(c *gin.Context) {
	var book model.Book
	
	// ShouldBindJSON автоматично перевіряє binding-теги, наприклад required для Title
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBook := handler.bookRepo.Create(book)
	c.JSON(http.StatusCreated, createdBook)
}