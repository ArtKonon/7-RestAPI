package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CategoryList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func CategoryGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
