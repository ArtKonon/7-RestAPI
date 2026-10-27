package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"homework/internal/handler"
	"homework/internal/middleware"
	"homework/internal/repository"
)

func NewRouter() http.Handler {
	router := gin.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	repo := repository.NewBookRepository()
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
