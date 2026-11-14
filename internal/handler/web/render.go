package web

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"libratrack/internal/model"
	"libratrack/internal/repository"
	pages "libratrack/web/templates/pages"
	"libratrack/web/templates/partials"

	"github.com/a-h/templ"
)

type Handler struct {
	repo repository.BookRepository
}

type Profile struct {
	Name     string
	Email    string
	Feedback string
}

const profileCookieName = "libratrack_profile"

func NewHandler(repo repository.BookRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) StartPage(w http.ResponseWriter, r *http.Request) {
	h.renderPage(w, r, pages.StartPage(), pages.StartContent(), http.StatusOK)
}

func (h *Handler) Clock(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format(time.RFC1123)
	if r.Header.Get("HX-Request") == "true" {
		h.renderPage(w, r, nil, pages.ClockBox(now), http.StatusOK)
		return
	}
	h.renderPage(w, r, pages.StartPage(), pages.StartContent(), http.StatusOK)
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	profile := readProfileCookie(r)
	h.renderPage(w, r, pages.LoginPage(profile.Name, profile.Email, profile.Feedback), pages.LoginContent(profile.Name, profile.Email, profile.Feedback), http.StatusOK)
}

func (h *Handler) SaveProfile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}
	profile := Profile{
		Name:     strings.TrimSpace(r.FormValue("name")),
		Email:    strings.TrimSpace(r.FormValue("email")),
		Feedback: strings.TrimSpace(r.FormValue("feedback")),
	}
	saveProfileCookie(w, profile)
	h.renderPage(w, r, pages.LoginPage(profile.Name, profile.Email, profile.Feedback), pages.LoginContent(profile.Name, profile.Email, profile.Feedback), http.StatusOK)
}

func saveProfileCookie(w http.ResponseWriter, profile Profile) {
	payload, err := json.Marshal(profile)
	if err != nil {
		return
	}
	encoded := base64.RawURLEncoding.EncodeToString(payload)
	http.SetCookie(w, &http.Cookie{
		Name:     profileCookieName,
		Value:    signCookie(encoded),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 30,
	})
}

func readProfileCookie(r *http.Request) Profile {
	cookie, err := r.Cookie(profileCookieName)
	if err != nil || cookie == nil {
		return Profile{}
	}
	encoded, ok := verifyCookie(cookie.Value)
	if !ok {
		return Profile{}
	}
	payload, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return Profile{}
	}
	var profile Profile
	if err := json.Unmarshal(payload, &profile); err != nil {
		return Profile{}
	}
	return profile
}

func signCookie(value string) string {
	mac := hmac.New(sha256.New, cookieSecret())
	mac.Write([]byte(value))
	return value + "." + hex.EncodeToString(mac.Sum(nil))
}

func verifyCookie(raw string) (string, bool) {
	parts := strings.Split(raw, ".")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", false
	}
	payload, signature := parts[0], parts[1]
	mac := hmac.New(sha256.New, cookieSecret())
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))
	return payload, hmac.Equal([]byte(signature), []byte(expected))
}

func cookieSecret() []byte {
	secret := strings.TrimSpace(os.Getenv("COOKIE_SECRET"))
	if secret == "" {
		secret = "libratrack-local-demo-secret"
	}
	return []byte(secret)
}

func (h *Handler) renderPage(w http.ResponseWriter, r *http.Request, full templ.Component, partial templ.Component, status int) {
	if r.Header.Get("HX-Request") == "true" && partial != nil {
		h.renderComponent(w, partial, status)
		return
	}
	if full == nil {
		full = partial
	}
	h.renderComponent(w, full, status)
}

func (h *Handler) renderComponent(w http.ResponseWriter, component templ.Component, status int) {
	if component == nil {
		w.WriteHeader(status)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := component.Render(context.Background(), w); err != nil {
		http.Error(w, "template render failed", http.StatusInternalServerError)
	}
}

func parseBookForm(r *http.Request) (partials.BookForm, error) {
	if err := r.ParseForm(); err != nil {
		return partials.BookForm{}, err
	}

	form := partials.BookForm{
		ID:       0,
		Title:    strings.TrimSpace(r.FormValue("title")),
		Author:   strings.TrimSpace(r.FormValue("author")),
		Category: strings.TrimSpace(r.FormValue("category")),
		Status:   strings.TrimSpace(r.FormValue("status")),
		Price:    strings.TrimSpace(r.FormValue("price")),
	}

	if form.Title == "" {
		return form, errors.New("title is required")
	}
	if form.Author == "" {
		return form, errors.New("author is required")
	}
	if form.Category == "" {
		return form, errors.New("category is required")
	}
	if form.Price != "" {
		if _, err := strconv.ParseFloat(form.Price, 64); err != nil {
			return form, errors.New("price must be a valid number")
		}
	}
	if form.Status == "" {
		form.Status = "available"
	}

	return form, nil
}

func toBookModel(form partials.BookForm) model.Book {
	status, price := normalizeBookStatus(form.Status), strings.TrimSpace(form.Price)
	book := model.Book{
		Title:    strings.TrimSpace(form.Title),
		ISBN:     normalizeISBN(form.Title),
		Authors:  []model.Author{{FullName: strings.TrimSpace(form.Author)}},
		Category: model.Category{Name: strings.TrimSpace(form.Category)},
	}
	if status != "" || price != "" {
		description := status + "|" + price
		book.Description = &description
	}
	return book
}

func toBookForm(book model.Book) partials.BookForm {
	status, price := parseStoredMetadata(book.Description)
	return partials.BookForm{
		ID:       book.ID,
		Title:    book.Title,
		Author:   book.LegacyAuthor(),
		Category: book.LegacyCategory(),
		Status:   status,
		Price:    price,
		Action:   "/books/" + strconv.Itoa(book.ID),
	}
}

func toBookItem(book model.Book) partials.BookItem {
	status, price := parseStoredMetadata(book.Description)
	return partials.BookItem{
		ID:       book.ID,
		Title:    book.Title,
		Author:   book.LegacyAuthor(),
		Category: book.LegacyCategory(),
		Status:   status,
		Price:    price,
	}
}

func normalizeISBN(title string) string {
	value := strings.TrimSpace(title)
	if value == "" {
		return "BOOK-0000"
	}
	value = strings.ToLower(value)
	value = strings.ReplaceAll(value, " ", "-")
	return "BOOK-" + value
}

func parseBookIDFromRequestPath(rawPath string) (int, error) {
	trimmed := strings.TrimSpace(rawPath)
	if trimmed == "" {
		return 0, errors.New("invalid book id")
	}
	segments := strings.Split(strings.Trim(trimmed, "/"), "/")
	if len(segments) == 0 {
		return 0, errors.New("invalid book id")
	}

	candidate := ""
	switch {
	case len(segments) >= 2 && segments[0] == "books":
		candidate = segments[1]
	case len(segments) >= 3 && segments[0] == "books" && segments[2] == "edit":
		candidate = segments[1]
	default:
		candidate = segments[len(segments)-1]
	}
	id, err := strconv.Atoi(candidate)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid book id")
	}
	return id, nil
}

func normalizeBookStatus(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "available":
		return "available"
	case "checked_out", "checked out":
		return "checked_out"
	case "archived":
		return "archived"
	default:
		return "available"
	}
}

func parseStoredMetadata(desc *string) (string, string) {
	status := "available"
	price := ""
	if desc == nil {
		return status, price
	}
	value := strings.TrimSpace(*desc)
	if value == "" {
		return status, price
	}
	parts := strings.SplitN(value, "|", 2)
	status = normalizeBookStatus(parts[0])
	if len(parts) == 2 {
		price = strings.TrimSpace(parts[1])
	}
	return status, price
}
