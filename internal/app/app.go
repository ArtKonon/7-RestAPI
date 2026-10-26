package app

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "homework/internal/handler"
    "homework/internal/repository"
)

func NewRouter() http.Handler {
    router := gin.Default()

    router.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if c.Request.Method == http.MethodOptions {
            c.Status(http.StatusNoContent)
            return
        }
        c.Next()
    })

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