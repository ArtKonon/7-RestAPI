package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"libratrack/internal/handler"
	"libratrack/internal/middleware"
	"libratrack/internal/repository"
)

func NewRouter() http.Handler {
	return NewRouterWithRepo(repository.NewInMemoryBookRepository())
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
	return NewRouterWithRepo(repository.NewGormBookRepository(db))
}
