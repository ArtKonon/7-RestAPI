package web

import (
	"net/http"
	"strconv"
	"strings"

	pages "libratrack/web/templates/pages"
	"libratrack/web/templates/partials"
)

func (h *Handler) ListBooksPage(w http.ResponseWriter, r *http.Request) {
	profile := readProfileCookie(r)
	books := h.repo.List()
	items := make([]partials.BookItem, 0, len(books))
	for _, book := range books {
		items = append(items, toBookItem(book))
	}
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	items = filterBookItems(items, query)
	page := pages.BooksPage(items, profile.Name, query)
	h.renderPage(w, r, page, pages.BooksContent(items, profile.Name, query), http.StatusOK)
}

func (h *Handler) NewBookPage(w http.ResponseWriter, r *http.Request) {
	form := partials.BookForm{
		Action: "/books",
		Status: "available",
	}
	page := pages.BookFormPage(form, "Create book")
	h.renderPage(w, r, page, pages.BookFormContent(form, "Create book"), http.StatusOK)
}

func (h *Handler) CreateBookPage(w http.ResponseWriter, r *http.Request) {
	h.NewBookPage(w, r)
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}
	form, err := parseBookForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if h.repo == nil {
		http.Error(w, "book repository unavailable", http.StatusInternalServerError)
		return
	}
	book := toBookModel(form)
	h.repo.Create(book)
	http.Redirect(w, r, "/books", http.StatusSeeOther)
}

func (h *Handler) EditBookPage(w http.ResponseWriter, r *http.Request) {
	id, err := parseBookIDFromRequestPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid book id", http.StatusBadRequest)
		return
	}
	if h.repo == nil {
		http.Error(w, "book repository unavailable", http.StatusInternalServerError)
		return
	}
	book, ok := h.repo.GetByID(id)
	if !ok {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}
	form := toBookForm(book)
	form.Action = "/books/" + strconv.Itoa(id)
	page := pages.BookFormPage(form, "Edit book")
	h.renderPage(w, r, page, pages.BookFormContent(form, "Edit book"), http.StatusOK)
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := parseBookIDFromRequestPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid book id", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}
	form, err := parseBookForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if h.repo == nil {
		http.Error(w, "book repository unavailable", http.StatusInternalServerError)
		return
	}
	book := toBookModel(form)
	if _, ok := h.repo.Update(id, book); !ok {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, "/books", http.StatusSeeOther)
}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := parseBookIDFromRequestPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid book id", http.StatusBadRequest)
		return
	}
	if h.repo == nil {
		http.Error(w, "book repository unavailable", http.StatusInternalServerError)
		return
	}
	if !h.repo.Delete(id) {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, "/books", http.StatusSeeOther)
}

func (h *Handler) SearchBooks(w http.ResponseWriter, r *http.Request) {
	term := strings.TrimSpace(r.URL.Query().Get("q"))
	books := h.repo.List()
	items := make([]partials.BookItem, 0, len(books))
	for _, book := range books {
		items = append(items, toBookItem(book))
	}
	filtered := filterBookItems(items, term)
	h.renderPage(w, r, pages.BooksPage(filtered, readProfileCookie(r).Name, term), pages.BooksContent(filtered, readProfileCookie(r).Name, term), http.StatusOK)
}

func filterBookItems(items []partials.BookItem, query string) []partials.BookItem {
	term := strings.TrimSpace(query)
	if term == "" {
		return items
	}
	needle := strings.ToLower(term)
	filtered := make([]partials.BookItem, 0, len(items))
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Title), needle) ||
			strings.Contains(strings.ToLower(item.Author), needle) ||
			strings.Contains(strings.ToLower(item.Category), needle) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
