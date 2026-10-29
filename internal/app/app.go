package app

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"libratrack/internal/handler"
	"libratrack/internal/middleware"
	"libratrack/internal/model"
	"libratrack/internal/repository"
)

func NewRouter() http.Handler {
	return NewRouterWithRepo(repository.NewInMemoryBookRepository())
}

func mountSwagger(router *gin.Engine) {
	wd, err := os.Getwd()
	if err != nil || wd == "" {
		wd = "."
	}

	swaggerDir := filepath.Join(wd, "swagger")
	router.Static("/swagger", swaggerDir)
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}

func NewRouterWithRepo(repo repository.BookRepository) http.Handler {
	router := gin.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	bookHandler := handler.NewBookHandler(repo)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	mountSwagger(router)

	api := router.Group("/api/v1")
	books := api.Group("/books")
	books.POST("", bookHandler.CreateBook)
	books.GET("", bookHandler.ListBooks)
	books.GET("/:id", bookHandler.GetBook)
	books.PUT("/:id", bookHandler.UpdateBook)
	books.DELETE("/:id", bookHandler.DeleteBook)

	return router
}

func NewRouterWithDB(db *gorm.DB) http.Handler {
	router := gin.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	mountSwagger(router)

	api := router.Group("/api/v1")
	registerCategoryRoutes(api, db)
	registerAuthorRoutes(api, db)
	registerBookRoutes(api, db)
	return router
}

func registerCategoryRoutes(api *gin.RouterGroup, db *gorm.DB) {
	if db == nil {
		return
	}
	categories := api.Group("/categories")

	categories.POST("", func(c *gin.Context) {
		var payload struct {
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			writeAPIError(c, http.StatusUnprocessableEntity, "validation_error", "invalid request body")
			return
		}
		name := strings.TrimSpace(payload.Name)
		if name == "" {
			writeAPIError(c, http.StatusUnprocessableEntity, "validation_error", "name is required")
			return
		}
		var existing model.Category
		if err := db.Where("LOWER(name) = LOWER(?)", name).First(&existing).Error; err == nil {
			writeAPIError(c, http.StatusConflict, "conflict", "category already exists")
			return
		}
		category := model.Category{Name: name}
		if err := db.Create(&category).Error; err != nil {
			writeAPIError(c, http.StatusInternalServerError, "internal_server_error", "internal server error")
			return
		}
		c.JSON(http.StatusCreated, category)
	})

	categories.GET("", func(c *gin.Context) {
		var categories []model.Category
		if err := db.Order("id").Find(&categories).Error; err != nil {
			writeAPIError(c, http.StatusInternalServerError, "internal_server_error", "internal server error")
			return
		}
		c.JSON(http.StatusOK, categories)
	})

	categories.GET("/:id", func(c *gin.Context) {
		id, err := parsePathID(c.Param("id"))
		if err != nil {
			writeAPIError(c, http.StatusBadRequest, "invalid_id", "invalid id")
			return
		}
		var category model.Category
		if err := db.First(&category, id).Error; err != nil {
			writeAPIError(c, http.StatusNotFound, "not_found", "category not found")
			return
		}
		c.JSON(http.StatusOK, category)
	})
}

func registerAuthorRoutes(api *gin.RouterGroup, db *gorm.DB) {
	if db == nil {
		return
	}
	authors := api.Group("/authors")

	authors.POST("", func(c *gin.Context) {
		var payload struct {
			FullName string `json:"full_name"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			writeAPIError(c, http.StatusUnprocessableEntity, "validation_error", "invalid request body")
			return
		}
		name := strings.TrimSpace(payload.FullName)
		if name == "" {
			writeAPIError(c, http.StatusUnprocessableEntity, "validation_error", "full_name is required")
			return
		}
		author := model.Author{FullName: name}
		if err := db.Create(&author).Error; err != nil {
			writeAPIError(c, http.StatusInternalServerError, "internal_server_error", "internal server error")
			return
		}
		c.JSON(http.StatusCreated, author)
	})

	authors.GET("", func(c *gin.Context) {
		var authors []model.Author
		if err := db.Order("id").Find(&authors).Error; err != nil {
			writeAPIError(c, http.StatusInternalServerError, "internal_server_error", "internal server error")
			return
		}
		c.JSON(http.StatusOK, authors)
	})

	authors.GET("/:id", func(c *gin.Context) {
		id, err := parsePathID(c.Param("id"))
		if err != nil {
			writeAPIError(c, http.StatusBadRequest, "invalid_id", "invalid id")
			return
		}
		var author model.Author
		if err := db.First(&author, id).Error; err != nil {
			writeAPIError(c, http.StatusNotFound, "not_found", "author not found")
			return
		}
		c.JSON(http.StatusOK, author)
	})
}

func registerBookRoutes(api *gin.RouterGroup, db *gorm.DB) {
	books := api.Group("/books")
	if db != nil {
		books.POST("", func(c *gin.Context) {
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
		})

		books.GET("", func(c *gin.Context) {
			page := 1
			limit := 10

			if q := c.Query("page"); q != "" {
				p, err := strconv.Atoi(q)
				if err != nil || p < 1 {
					writeAPIError(c, http.StatusBadRequest, "invalid_pagination", "invalid pagination")
					return
				}
				page = p
			}

			if q := c.Query("limit"); q != "" {
				l, err := strconv.Atoi(q)
				if err != nil || l < 1 {
					writeAPIError(c, http.StatusBadRequest, "invalid_pagination", "invalid pagination")
					return
				}
				if l > 100 {
					l = 100
				}
				limit = l
			}

			var books []model.Book
			query := db.Preload("Authors").Preload("Category").Order("id")
			if categoryID := strings.TrimSpace(c.Query("category_id")); categoryID != "" {
				if id, err := strconv.Atoi(categoryID); err == nil && id > 0 {
					query = query.Where("category_id = ?", id)
				}
			}
			if err := query.Find(&books).Error; err != nil {
				writeAPIError(c, http.StatusInternalServerError, "internal_server_error", "internal server error")
				return
			}

			start := (page - 1) * limit
			if start >= len(books) {
				c.JSON(http.StatusOK, []gin.H{})
				return
			}

			end := start + limit
			if end > len(books) {
				end = len(books)
			}

			resp := make([]gin.H, 0, end-start)
			for _, book := range books[start:end] {
				resp = append(resp, buildBookResponse(book))
			}
			c.JSON(http.StatusOK, resp)
		})

		books.GET("/:id", func(c *gin.Context) {
			id, err := parsePathID(c.Param("id"))
			if err != nil {
				writeAPIError(c, http.StatusBadRequest, "invalid_id", "invalid id")
				return
			}
			var book model.Book
			if err := db.Preload("Authors").Preload("Category").First(&book, id).Error; err != nil {
				writeAPIError(c, http.StatusNotFound, "not_found", "book not found")
				return
			}
			c.JSON(http.StatusOK, buildBookResponse(book))
		})

		books.PUT("/:id", func(c *gin.Context) {
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
		})

		books.DELETE("/:id", func(c *gin.Context) {
			id, err := parsePathID(c.Param("id"))
			if err != nil {
				writeAPIError(c, http.StatusBadRequest, "invalid_id", "invalid id")
				return
			}
			if db.Delete(&model.Book{}, id).RowsAffected == 0 {
				writeAPIError(c, http.StatusNotFound, "not_found", "book not found")
				return
			}
			c.Status(http.StatusNoContent)
		})

		books.POST("/:id/copies", func(c *gin.Context) {
			id, err := parsePathID(c.Param("id"))
			if err != nil {
				writeAPIError(c, http.StatusBadRequest, "invalid_id", "invalid id")
				return
			}
			var book model.Book
			if err := db.First(&book, id).Error; err != nil {
				writeAPIError(c, http.StatusNotFound, "not_found", "book not found")
				return
			}
			var payload struct {
				InventoryNo string `json:"inventory_no"`
				Status      string `json:"status"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				writeAPIError(c, http.StatusUnprocessableEntity, "validation_error", "invalid request body")
				return
			}
			inventoryNo := strings.TrimSpace(payload.InventoryNo)
			status := strings.TrimSpace(payload.Status)
			if inventoryNo == "" {
				writeAPIError(c, http.StatusUnprocessableEntity, "validation_error", "inventory_no is required")
				return
			}
			if status == "" {
				status = "available"
			}
			if status != "available" && status != "lost" && status != "damaged" {
				writeAPIError(c, http.StatusUnprocessableEntity, "validation_error", "invalid copy status")
				return
			}
			var existing model.Copy
			if err := db.Where("inventory_no = ?", inventoryNo).First(&existing).Error; err == nil {
				writeAPIError(c, http.StatusConflict, "conflict", "inventory_no already exists")
				return
			}
			copy := model.Copy{BookID: book.ID, InventoryNo: inventoryNo, Status: status}
			if err := db.Create(&copy).Error; err != nil {
				writeAPIError(c, http.StatusInternalServerError, "internal_server_error", "internal server error")
				return
			}
			c.JSON(http.StatusCreated, copy)
		})

		books.GET("/:id/copies", func(c *gin.Context) {
			id, err := parsePathID(c.Param("id"))
			if err != nil {
				writeAPIError(c, http.StatusBadRequest, "invalid_id", "invalid id")
				return
			}
			var book model.Book
			if err := db.First(&book, id).Error; err != nil {
				writeAPIError(c, http.StatusNotFound, "not_found", "book not found")
				return
			}
			var copies []model.Copy
			if err := db.Where("book_id = ?", book.ID).Order("id").Find(&copies).Error; err != nil {
				writeAPIError(c, http.StatusInternalServerError, "internal_server_error", "internal server error")
				return
			}
			c.JSON(http.StatusOK, copies)
		})

		books.GET("/:id/availability", func(c *gin.Context) {
			id, err := parsePathID(c.Param("id"))
			if err != nil {
				writeAPIError(c, http.StatusBadRequest, "invalid_id", "invalid id")
				return
			}
			var book model.Book
			if err := db.First(&book, id).Error; err != nil {
				writeAPIError(c, http.StatusNotFound, "not_found", "book not found")
				return
			}
			var total int64
			if err := db.Model(&model.Copy{}).Where("book_id = ?", book.ID).Count(&total).Error; err != nil {
				writeAPIError(c, http.StatusInternalServerError, "internal_server_error", "internal server error")
				return
			}
			var available int64
			if err := db.Model(&model.Copy{}).Where("book_id = ? AND status = ?", book.ID, "available").Count(&available).Error; err != nil {
				writeAPIError(c, http.StatusInternalServerError, "internal_server_error", "internal server error")
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"book_id":          book.ID,
				"total_copies":     total,
				"available_copies": available,
			})
		})
	}
}

func buildBookResponse(book model.Book) gin.H {
	category := gin.H{"id": book.Category.ID, "name": book.Category.Name}
	authors := make([]gin.H, 0, len(book.Authors))
	for _, author := range book.Authors {
		authors = append(authors, gin.H{"id": author.ID, "full_name": author.FullName})
	}
	response := gin.H{
		"id":             book.ID,
		"title":          book.Title,
		"isbn":           book.ISBN,
		"category_id":    book.CategoryID,
		"category":       category,
		"authors":        authors,
		"published_year": book.PublishedYear,
		"description":    book.Description,
		"created_at":     book.CreatedAt.Format(time.RFC3339Nano),
		"updated_at":     book.UpdatedAt.Format(time.RFC3339Nano),
	}
	if book.PublishedYear == nil {
		response["published_year"] = nil
	}
	if book.Description == nil {
		response["description"] = nil
	}
	return response
}

func parsePathID(raw string) (int, error) {
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func writeAPIError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}
