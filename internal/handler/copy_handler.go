package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CopyList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func CopyCreate(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"status": "ok"})
}
