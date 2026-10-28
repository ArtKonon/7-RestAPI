package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func AuthorGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
