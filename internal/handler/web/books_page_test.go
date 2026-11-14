package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"libratrack/internal/model"
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

func TestBooksPageFiltersByQueryAndPreservesSearchValue(t *testing.T) {
	repo := repository.NewInMemoryBookRepository()
	repo.Create(model.Book{Title: "Clean Architecture", ISBN: "1", Category: model.Category{Name: "Design"}, Authors: []model.Author{{FullName: "Robert Martin"}}})
	repo.Create(model.Book{Title: "Go in Practice", ISBN: "2", Category: model.Category{Name: "Programming"}, Authors: []model.Author{{FullName: "Matt Butcher"}}})

	h := NewHandler(repo)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/books?q=go", nil)

	h.ListBooksPage(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	body := res.Body.String()
	if !strings.Contains(body, "Go in Practice") {
		t.Fatalf("expected filtered book in page, got: %s", body)
	}
	if strings.Contains(body, "Clean Architecture") {
		t.Fatalf("expected non-matching book to be filtered out, got: %s", body)
	}
	if !strings.Contains(body, "value=\"go\"") {
		t.Fatalf("expected books search query preserved in form, got: %s", body)
	}
}
