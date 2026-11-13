package app

import (
	"net/http"
	"strings"
	"time"

	"libratrack/internal/handler"
	"libratrack/internal/middleware"
	"libratrack/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerAuthRoutes(api *gin.RouterGroup, db *gorm.DB) {
	if db == nil {
		return
	}

	public := api.Group("")
	public.POST("/auth/register", handler.RegisterUser(db))
	public.POST("/auth/login", handler.LoginUser(db))
}

func registerProtectedRoutes(api *gin.RouterGroup, db *gorm.DB) {
	if db == nil {
		return
	}

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/me/loans", func(c *gin.Context) {
		handler.MeLoans(c, db)
	})
}

func registerStaffRoutes(api *gin.RouterGroup, db *gorm.DB) {
	if db == nil {
		return
	}

	staff := api.Group("")
	staff.Use(middleware.AuthMiddleware(), middleware.RequireRole("librarian", "admin"))
	staff.POST("/loans", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"status": "loan_created"})
	})
	staff.POST("/loans/:id/return", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "loan_returned"})
	})
}

func registerAdminRoutes(api *gin.RouterGroup, db *gorm.DB) {
	if db == nil {
		return
	}

	admin := api.Group("")
	admin.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	admin.GET("/users", handler.ListUsers(db))
	admin.PUT("/users/:id/role", handler.UpdateUserRole(db))
}

func createBookRestrictedHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			Title         string  `json:"title"`
			ISBN          string  `json:"isbn"`
			CategoryID    int     `json:"category_id"`
			AuthorIDs     []int   `json:"author_ids"`
			PublishedYear *int    `json:"published_year"`
			Description   *string `json:"description"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			writeAPIError(c, http.StatusUnprocessableEntity, "validation_error", "invalid request body")
			return
		}
		if strings.TrimSpace(payload.Title) == "" || strings.TrimSpace(payload.ISBN) == "" || payload.CategoryID <= 0 || len(payload.AuthorIDs) == 0 {
			writeAPIError(c, http.StatusUnprocessableEntity, "validation_error", "missing required fields")
			return
		}
		var category model.Category
		if err := db.First(&category, payload.CategoryID).Error; err != nil {
			writeAPIError(c, http.StatusUnprocessableEntity, "invalid_reference", "category not found")
			return
		}
		authors := make([]model.Author, 0, len(payload.AuthorIDs))
		for _, authorID := range payload.AuthorIDs {
			var author model.Author
			if err := db.First(&author, authorID).Error; err != nil {
				writeAPIError(c, http.StatusUnprocessableEntity, "invalid_reference", "author not found")
				return
			}
			authors = append(authors, author)
		}
		var existing model.Book
		if err := db.Where("isbn = ?", payload.ISBN).First(&existing).Error; err == nil {
			writeAPIError(c, http.StatusConflict, "conflict", "book already exists")
			return
		}
		book := model.Book{
			Title:         payload.Title,
			ISBN:          payload.ISBN,
			CategoryID:    payload.CategoryID,
			Category:      category,
			Authors:       authors,
			PublishedYear: payload.PublishedYear,
			Description:   payload.Description,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := db.Create(&book).Error; err != nil {
			writeAPIError(c, http.StatusInternalServerError, "internal_server_error", "internal server error")
			return
		}
		if err := db.Model(&book).Association("Authors").Replace(authors); err != nil {
			writeAPIError(c, http.StatusInternalServerError, "internal_server_error", "internal server error")
			return
		}
		var saved model.Book
		if err := db.Preload("Authors").Preload("Category").First(&saved, book.ID).Error; err != nil {
			writeAPIError(c, http.StatusInternalServerError, "internal_server_error", "internal server error")
			return
		}
		c.JSON(http.StatusCreated, buildBookResponse(saved))
	}
}

func updateBookRestrictedHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := parsePathID(c.Param("id"))
		if err != nil {
			writeAPIError(c, http.StatusBadRequest, "invalid_id", "invalid id")
			return
		}
		var existing model.Book
		if err := db.First(&existing, id).Error; err != nil {
			writeAPIError(c, http.StatusNotFound, "not_found", "book not found")
			return
		}
		var payload struct {
			Title         string  `json:"title"`
			ISBN          string  `json:"isbn"`
			CategoryID    int     `json:"category_id"`
			AuthorIDs     []int   `json:"author_ids"`
			PublishedYear *int    `json:"published_year"`
			Description   *string `json:"description"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			writeAPIError(c, http.StatusUnprocessableEntity, "validation_error", "invalid request body")
			return
		}
		if strings.TrimSpace(payload.Title) == "" || strings.TrimSpace(payload.ISBN) == "" || payload.CategoryID <= 0 || len(payload.AuthorIDs) == 0 {
			writeAPIError(c, http.StatusUnprocessableEntity, "validation_error", "missing required fields")
			return
		}
		var category model.Category
		if err := db.First(&category, payload.CategoryID).Error; err != nil {
			writeAPIError(c, http.StatusUnprocessableEntity, "invalid_reference", "category not found")
			return
		}
		authors := make([]model.Author, 0, len(payload.AuthorIDs))
		for _, authorID := range payload.AuthorIDs {
			var author model.Author
			if err := db.First(&author, authorID).Error; err != nil {
				writeAPIError(c, http.StatusUnprocessableEntity, "invalid_reference", "author not found")
				return
			}
			authors = append(authors, author)
		}
		existing.Title = payload.Title
		existing.ISBN = payload.ISBN
		existing.CategoryID = payload.CategoryID
		existing.Category = category
		existing.Authors = authors
		existing.PublishedYear = payload.PublishedYear
		existing.Description = payload.Description
		existing.UpdatedAt = time.Now().UTC()
		if err := db.Save(&existing).Error; err != nil {
			writeAPIError(c, http.StatusInternalServerError, "internal_server_error", "internal server error")
			return
		}
		if err := db.Model(&existing).Association("Authors").Replace(authors); err != nil {
			writeAPIError(c, http.StatusInternalServerError, "internal_server_error", "internal server error")
			return
		}
		var saved model.Book
		if err := db.Preload("Authors").Preload("Category").First(&saved, id).Error; err != nil {
			writeAPIError(c, http.StatusInternalServerError, "internal_server_error", "internal server error")
			return
		}
		c.JSON(http.StatusOK, buildBookResponse(saved))
	}
}
