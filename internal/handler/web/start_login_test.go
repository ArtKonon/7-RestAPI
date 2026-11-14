package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"libratrack/internal/model"
	"libratrack/internal/repository"
)

func TestStartPageIncludesHTMXIntroAndClock(t *testing.T) {
	h := NewHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	h.StartPage(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	body := res.Body.String()
	if !strings.Contains(body, "HTMX") {
		t.Fatal("expected start page to mention HTMX")
	}
	if !strings.Contains(body, "every 2s") {
		t.Fatal("expected start page to include a 2-second htmx clock")
	}
}

func TestLoginPageSavesProfileAndFeedback(t *testing.T) {
	h := NewHandler(nil)
	profile := Profile{Name: "Jane Doe", Email: "jane@example.com", Feedback: "Need more classics"}
	res := httptest.NewRecorder()
	ssr := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("name=Jane+Doe&email=jane%40example.com&feedback=Need+more+classics"))
	ssr.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	h.SaveProfile(res, ssr)
	if res.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect 303, got %d", res.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	for _, c := range res.Result().Cookies() {
		req.AddCookie(c)
	}
	got := readProfileCookie(req)
	if got.Name != profile.Name || got.Email != profile.Email || got.Feedback != profile.Feedback {
		t.Fatalf("profile mismatch: got %#v want %#v", got, profile)
	}
}

func TestAuthorsPageUsesStoredBioAndSearch(t *testing.T) {
	repo := repository.NewInMemoryBookRepository()
	bio := "This is the actual author bio from the database."
	repo.Create(model.Book{
		Title:    "The Hound of the Baskervilles",
		ISBN:     "978-1-2345-6789-0",
		Category: model.Category{Name: "Classics"},
		Authors:  []model.Author{{FullName: "Arthur Conan Doyle", Bio: &bio}},
	})

	h := NewHandler(repo)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/authors?q=conan", nil)

	h.AuthorsPage(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	body := res.Body.String()
	if !strings.Contains(body, "This is the actual author bio from the database.") {
		t.Fatalf("expected stored author bio in page, got %s", body)
	}
	if !strings.Contains(body, "Arthur Conan Doyle") {
		t.Fatalf("expected author name in page, got %s", body)
	}
	if !strings.Contains(body, "value=\"conan\"") {
		t.Fatalf("expected authors search query preserved in form, got %s", body)
	}
}
