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
	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	body := res.Body.String()
	if !strings.Contains(body, "Jane Doe") || !strings.Contains(body, "jane@example.com") {
		t.Fatalf("expected saved profile summary to render immediately, got %s", body)
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

func TestSaveProfileRendersSummaryImmediately(t *testing.T) {
	h := NewHandler(repository.NewInMemoryBookRepository())
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("name=Jane+Doe&email=jane%40example.com&feedback=Need+more+classics"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	h.SaveProfile(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	body := res.Body.String()
	if !strings.Contains(body, "Saved profile: Jane Doe · jane@example.com") {
		t.Fatalf("expected inline profile summary after save, got: %s", body)
	}
	if strings.Contains(body, "value=\"Jane Doe\"") || strings.Contains(body, "value=\"jane@example.com\"") {
		t.Fatalf("expected form values to clear after save, got: %s", body)
	}
}

func TestUpdateAuthorPropagatesToBooks(t *testing.T) {
	repo := repository.NewInMemoryBookRepository()
	bio := "Original bio."
	book := model.Book{
		Title:    "The Hound of the Baskervilles",
		ISBN:     "978-1-2345-6789-0",
		Category: model.Category{Name: "Classics"},
		Authors:  []model.Author{{FullName: "Arthur Conan Doyle", Bio: &bio}},
	}
	repo.Create(book)

	h := NewHandler(repo)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/authors/Arthur+Conan+Doyle", strings.NewReader("old_name=Arthur+Conan+Doyle&name=Sir+Arthur+Conan+Doyle&bio=Updated+bio."))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawPath = "/authors/Arthur+Conan+Doyle"

	h.UpdateAuthor(res, req)
	if res.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect 303, got %d", res.Code)
	}
	updated, ok := repo.GetByID(1)
	if !ok {
		t.Fatal("expected book to exist")
	}
	if len(updated.Authors) != 1 || updated.Authors[0].FullName != "Sir Arthur Conan Doyle" {
		t.Fatalf("expected author name updated in book, got %#v", updated.Authors)
	}
	if updated.Authors[0].Bio == nil || *updated.Authors[0].Bio != "Updated bio." {
		t.Fatalf("expected author bio updated in book, got %#v", updated.Authors[0].Bio)
	}
}

func TestAuthorsPageDoesNotExposeCreateActionsOrPublishedWarning(t *testing.T) {
	repo := repository.NewInMemoryBookRepository()
	book := model.Book{
		Title:    "The Hound of the Baskervilles",
		ISBN:     "978-1-2345-6789-0",
		Category: model.Category{Name: "Classics"},
		Authors:  []model.Author{{FullName: "Arthur Conan Doyle", Bio: stringPtr("A classic author.")}},
	}
	repo.Create(book)

	h := NewHandler(repo)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/authors", nil)
	res = httptest.NewRecorder()
	if err := req.ParseForm(); err != nil {
		t.Fatal(err)
	}

	h.AuthorsPage(res, req)
	body := res.Body.String()
	if strings.Contains(body, "Create author") {
		t.Fatalf("authors page should not expose create author action, got: %s", body)
	}
	if strings.Contains(body, "This author is published.") {
		t.Fatalf("authors page should not expose published warning, got: %s", body)
	}
	if strings.Contains(body, "/books/") && strings.Contains(body, "Edit") {
		t.Fatalf("authors page should not expose book edit links, got: %s", body)
	}
}

func stringPtr(s string) *string {
	return &s
}
