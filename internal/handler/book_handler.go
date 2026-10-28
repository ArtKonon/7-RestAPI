package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"libratrack/internal/model"
	"libratrack/internal/repository"
)

type BookHandler struct {
	bookRepo repository.BookRepository
}

func NewBookHandler(repo repository.BookRepository) *BookHandler {
	return &BookHandler{bookRepo: repo}
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var book model.Book

	err := c.ShouldBindJSON(&book)
	if err != nil {
		writeError(c, http.StatusUnprocessableEntity, "validation_error", "invalid request body")
		return
	}

	if strings.TrimSpace(book.Title) == "" ||
		strings.TrimSpace(book.ISBN) == "" ||
		len(book.Authors) == 0 && strings.TrimSpace(book.LegacyAuthor()) == "" ||
		book.CategoryID <= 0 && strings.TrimSpace(book.LegacyCategory()) == "" {
		writeError(c, http.StatusUnprocessableEntity, "validation_error", "missing required fields")
		return
	}

	created := h.bookRepo.Create(book)
	c.JSON(http.StatusCreated, created)
}

func (h *BookHandler) GetBook(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid_id", "invalid id")
		return
	}

	book, ok := h.bookRepo.GetByID(id)
	if !ok {
		writeError(c, http.StatusNotFound, "not_found", "Book not found")
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) ListBooks(c *gin.Context) {
	// page и limit как в README
	page := 1
	limit := 10

	if q := c.Query("page"); q != "" {
		p, err := strconv.Atoi(q)
		if err != nil || p < 1 {
			writeError(c, http.StatusBadRequest, "invalid_pagination", "invalid pagination")
			return
		}
		page = p
	}

	if q := c.Query("limit"); q != "" {
		l, err := strconv.Atoi(q)
		if err != nil || l < 1 {
			writeError(c, http.StatusBadRequest, "invalid_pagination", "invalid pagination")
			return
		}
		if l > 100 {
			l = 100
		}
		limit = l
	}

	books := h.bookRepo.List()

	if filter := strings.TrimSpace(c.Query("category")); filter != "" {
		filtered := make([]model.Book, 0, len(books))
		for _, book := range books {
			if strings.EqualFold(book.Category.Name, filter) || strings.EqualFold(book.LegacyCategory(), filter) {
				filtered = append(filtered, book)
			}
		}
		books = filtered
	}

	start := (page - 1) * limit
	if start > len(books) {
		c.JSON(http.StatusOK, []model.Book{})
		return
	}

	end := start + limit
	if end > len(books) {
		end = len(books)
	}

	c.JSON(http.StatusOK, books[start:end])
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid_id", "invalid id")
		return
	}

	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		writeError(c, http.StatusUnprocessableEntity, "validation_error", "invalid request body")
		return
	}

	if strings.TrimSpace(book.Title) == "" ||
		strings.TrimSpace(book.ISBN) == "" ||
		len(book.Authors) == 0 && strings.TrimSpace(book.LegacyAuthor()) == "" ||
		book.CategoryID <= 0 && strings.TrimSpace(book.LegacyCategory()) == "" {
		writeError(c, http.StatusUnprocessableEntity, "validation_error", "missing required fields")
		return
	}

	updated, ok := h.bookRepo.Update(id, book)
	if !ok {
		writeError(c, http.StatusNotFound, "not_found", "Book not found")
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid_id", "invalid id")
		return
	}

	if !h.bookRepo.Delete(id) {
		writeError(c, http.StatusNotFound, "not_found", "Book not found")
		return
	}

	c.Status(http.StatusNoContent)
}

func parseID(s string) (int, error) {
	id, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, errors.New("invalid id")
	}
	if id > uint64(^uint32(0)) {
		return 0, errors.New("invalid id")
	}
	return int(id), nil
}

func writeError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}
