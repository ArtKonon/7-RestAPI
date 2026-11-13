package handler

import (
	"net/http"
	"strings"
	"time"

	"libratrack/internal/auth"
	"libratrack/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"error": gin.H{"code": "validation_error", "message": "invalid request body"},
			})
			return
		}

		email := strings.TrimSpace(strings.ToLower(req.Email))
		password := strings.TrimSpace(req.Password)
		if email == "" || password == "" {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"error": gin.H{"code": "validation_error", "message": "email and password are required"},
			})
			return
		}

		var existing model.User
		if err := db.Where("LOWER(email) = LOWER(?)", email).First(&existing).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"error": gin.H{"code": "duplicate_email", "message": "email already registered"},
			})
			return
		}

		hash, err := auth.HashPassword(password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "internal_server_error", "message": "internal server error"},
			})
			return
		}

		user := model.User{
			Email:        email,
			PasswordHash: hash,
			Role:         "reader",
		}
		if err := db.Create(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "internal_server_error", "message": "internal server error"},
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		})
	}
}

func LoginUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{"code": "unauthorized", "message": "invalid credentials"},
			})
			return
		}

		email := strings.TrimSpace(strings.ToLower(req.Email))
		password := strings.TrimSpace(req.Password)
		if email == "" || password == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{"code": "unauthorized", "message": "invalid credentials"},
			})
			return
		}

		var user model.User
		if err := db.Where("LOWER(email) = LOWER(?)", email).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{"code": "unauthorized", "message": "invalid credentials"},
			})
			return
		}

		if err := auth.ComparePassword(password, user.PasswordHash); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{"code": "unauthorized", "message": "invalid credentials"},
			})
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Role, 24*time.Hour)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "internal_server_error", "message": "internal server error"},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
				"role":  user.Role,
			},
		})
	}
}

func MeLoans(c *gin.Context, db *gorm.DB) {
	_ = db
	userID, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{"code": "unauthorized", "message": "unauthorized"},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"loans":   []gin.H{},
	})
}

func ListUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []model.User
		if err := db.Order("id").Find(&users).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "internal_server_error", "message": "internal server error"},
			})
			return
		}
		out := make([]gin.H, 0, len(users))
		for _, user := range users {
			out = append(out, gin.H{
				"id":    user.ID,
				"email": user.Email,
				"role":  user.Role,
			})
		}
		c.JSON(http.StatusOK, out)
	}
}

func UpdateUserRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var user model.User
		if err := db.First(&user, id).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "not_found", "message": "user not found"},
			})
			return
		}

		var payload struct {
			Role string `json:"role"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"error": gin.H{"code": "validation_error", "message": "invalid request body"},
			})
			return
		}

		role := strings.TrimSpace(strings.ToLower(payload.Role))
		if role != "reader" && role != "librarian" && role != "admin" {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"error": gin.H{"code": "validation_error", "message": "role must be one of reader, librarian, admin"},
			})
			return
		}

		user.Role = role
		if err := db.Save(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "internal_server_error", "message": "internal server error"},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		})
	}
}
