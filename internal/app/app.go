// Package app wires the whole HTTP API together.
//
// The automatic tests talk to your API ONLY through NewRouter(), so you are free to
// organise the rest of the code (model, repository, handler, middleware) as you like.
package app

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// NewRouter must return a fully configured HTTP handler for your resource.
//
// Requirements (see README.md for the full contract):
//   - every call returns a NEW router with its own EMPTY in-memory storage;
//   - it must be safe for concurrent requests;
//   - it can be built with net/http, gin, chi or any other router.
//



func NewRouter() http.Handler {

	router := gin.Default() // підключає middleware: Logger і Recovery

	books := router.Group("/books")
	{ 
		books.GET("", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
		})
		// GET /books/{id} - отримати книгу за ID
		books.GET("/:id", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
		})
		
		// POST /books - додати нову книгу
		books.POST("", handler.CreateBook)
		
		// PUT /books/{id} - оновити книгу
		books.PUT("/:id", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
		})
		
		// DELETE /books/{id} - видалити книгу
		books.DELETE("/:id", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
		})
	}

	return router


}
