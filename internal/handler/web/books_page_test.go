package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"libratrack/internal/repository"
)

func TestBooksPageContainsMainNavigation(t *testing.T) {
	h := NewHandler(repository.NewInMemoryBookRepository())
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/books", nil)

	h.ListBooksPage(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	body := res.Body.String()
	if !strings.Contains(body, "href=\"/\"") || !strings.Contains(body, "href=\"/authors\"") || !strings.Contains(body, "href=\"/login\"") {
		t.Fatalf("expected top navigation links for Start, Authors and Login in books page, got: %s", body)
	}
}

func TestCreateBookPageUsesLeftAlignedBackLink(t *testing.T) {
	h := NewHandler(repository.NewInMemoryBookRepository())
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/books/new", nil)

	h.NewBookPage(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	body := res.Body.String()
	if !strings.Contains(body, "Back to list") {
		t.Fatalf("expected Back to list link in create form, got: %s", body)
	}
	if !strings.Contains(body, "form-actions left") {
		t.Fatalf("expected left aligned form actions in create page, got: %s", body)
	}
}
