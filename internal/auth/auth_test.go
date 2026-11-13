package auth_test

import (
	"os"
	"testing"
	"time"

	"libratrack/internal/auth"
)

func TestHashPasswordAndCompare(t *testing.T) {
	tests := []struct {
		name       string
		password   string
		wrongPass  string
		wantHashOK bool
		wantCmpOK  bool
	}{
		{name: "valid hashing and comparison", password: "super-secret", wrongPass: "wrong-pass", wantHashOK: true, wantCmpOK: true},
		{name: "empty password still hashed", password: "", wrongPass: "x", wantHashOK: true, wantCmpOK: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := auth.HashPassword(tt.password)
			if err != nil {
				t.Fatalf("HashPassword returned error: %v", err)
			}
			if tt.wantHashOK && hash == tt.password {
				t.Fatalf("HashPassword should not return the raw password")
			}
			if err := auth.ComparePassword(tt.password, hash); err != nil {
				t.Fatalf("ComparePassword should accept the same password: %v", err)
			}
			if err := auth.ComparePassword(tt.wrongPass, hash); err == nil && tt.wrongPass != "" {
				t.Fatal("ComparePassword should reject a wrong password")
			}
		})
	}
}

func TestGenerateAndParseTokenTable(t *testing.T) {
	if err := os.Setenv("JWT_SECRET", "test-secret"); err != nil {
		t.Fatalf("Setenv failed: %v", err)
	}
	defer os.Unsetenv("JWT_SECRET")

	tests := []struct {
		name     string
		userID   int
		role     string
		ttl      time.Duration
		wantRole string
		wantErr  bool
	}{
		{name: "admin token", userID: 42, role: "admin", ttl: 24 * time.Hour, wantRole: "admin"},
		{name: "reader token with custom ttl", userID: 7, role: "reader", ttl: 2 * time.Hour, wantRole: "reader"},
		{name: "zero ttl uses default", userID: 1, role: "librarian", ttl: 0, wantRole: "librarian"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := auth.GenerateToken(tt.userID, tt.role, tt.ttl)
			if err != nil {
				t.Fatalf("GenerateToken returned error: %v", err)
			}

			claims, err := auth.ParseToken(token)
			if err != nil {
				t.Fatalf("ParseToken returned error: %v", err)
			}
			if claims.UserID != tt.userID {
				t.Fatalf("expected user_id %d, got %d", tt.userID, claims.UserID)
			}
			if claims.Role != tt.wantRole {
				t.Fatalf("expected role %q, got %q", tt.wantRole, claims.Role)
			}
			if claims.ExpiresAt == nil {
				t.Fatal("token expiry should be set")
			}
		})
	}
}
