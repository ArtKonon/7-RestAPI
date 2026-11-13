package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"libratrack/internal/middleware"

	"github.com/gin-gonic/gin"
)

func TestRequireRoleTable(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		role        string
		allowed     []string
		wantStatus  int
		wantMessage string
	}{
		{name: "librarian allowed", role: "librarian", allowed: []string{"librarian", "admin"}, wantStatus: http.StatusOK},
		{name: "admin allowed", role: "admin", allowed: []string{"librarian", "admin"}, wantStatus: http.StatusOK},
		{name: "reader denied", role: "reader", allowed: []string{"librarian", "admin"}, wantStatus: http.StatusForbidden},
		{name: "missing role denied", role: "", allowed: []string{"admin"}, wantStatus: http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			if tt.role != "" {
				c.Set("role", tt.role)
			}
			c.Request = httptest.NewRequest(http.MethodPost, "/books", nil)

			middleware.RequireRole(tt.allowed...)(c)

			if w.Code != tt.wantStatus {
				t.Fatalf("role=%q allowed=%v: expected status %d, got %d", tt.role, tt.allowed, tt.wantStatus, w.Code)
			}
		})
	}
}
