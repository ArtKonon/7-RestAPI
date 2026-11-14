package web

import (
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"libratrack/internal/model"
	pages "libratrack/web/templates/pages"
)

func (h *Handler) AuthorsPage(w http.ResponseWriter, r *http.Request) {
	profile := readProfileCookie(r)
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	alert := ""
	if r.URL.Query().Get("error") == "published" {
		alert = "This author is published."
	}
	cards := collectAuthorCards(h.repo.List())
	filtered := make([]pages.AuthorCard, 0, len(cards))
	for _, card := range cards {
		if query != "" && !matchesAuthorQuery(card.Name, card.Books, card.Bio, query) {
			continue
		}
		filtered = append(filtered, card)
	}
	page := pages.AuthorsPage(filtered, profile.Name, query, alert)
	h.renderPage(w, r, page, pages.AuthorsContent(filtered, profile.Name, query, alert), http.StatusOK)
}

func (h *Handler) NewAuthorPage(w http.ResponseWriter, r *http.Request) {
	profile := readProfileCookie(r)
	page := pages.AuthorFormPage(pages.AuthorForm{UserName: profile.Name}, "Create author")
	h.renderPage(w, r, page, pages.AuthorFormContent(pages.AuthorForm{UserName: profile.Name}, "Create author"), http.StatusOK)
}

func (h *Handler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	bio := strings.TrimSpace(r.FormValue("bio"))
	if name == "" {
		http.Error(w, "author name is required", http.StatusBadRequest)
		return
	}
	book := model.Book{
		Title:    "Author profile: " + name,
		ISBN:     "AUTHOR-" + strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(name), " ", "-")),
		Category: model.Category{Name: "Author"},
		Authors:  []model.Author{{FullName: name, Bio: authorBioValue(bio)}},
	}
	h.repo.Create(book)
	http.Redirect(w, r, "/authors", http.StatusSeeOther)
}

func (h *Handler) EditAuthorPage(w http.ResponseWriter, r *http.Request) {
	profile := readProfileCookie(r)
	name := strings.TrimSpace(parseAuthorNameFromPath(r.URL.Path, "edit"))
	cards := collectAuthorCards(h.repo.List())
	for _, card := range cards {
		if card.Name == name {
			page := pages.AuthorEditPage(card, profile.Name)
			h.renderPage(w, r, page, pages.AuthorEditContent(card, profile.Name), http.StatusOK)
			return
		}
	}
	http.NotFound(w, r)
}

func (h *Handler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}
	oldName := strings.TrimSpace(r.FormValue("old_name"))
	if oldName == "" {
		oldName = strings.TrimSpace(parseAuthorNameFromPath(r.URL.Path, ""))
	}
	newName := strings.TrimSpace(r.FormValue("name"))
	bio := strings.TrimSpace(r.FormValue("bio"))
	if oldName == "" || newName == "" {
		http.Error(w, "author name is required", http.StatusBadRequest)
		return
	}
	books := h.repo.List()
	updated := false
	for _, book := range books {
		changed := false
		for i, author := range book.Authors {
			if !strings.EqualFold(strings.TrimSpace(author.FullName), oldName) {
				continue
			}
			author.FullName = newName
			author.Bio = authorBioValue(bio)
			book.Authors[i] = author
			changed = true
		}
		if changed {
			h.repo.Update(book.ID, book)
			updated = true
		}
	}
	if !updated {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/authors", http.StatusSeeOther)
}

func (h *Handler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(parseAuthorNameFromPath(r.URL.Path, "delete"))
	if name == "" {
		http.Error(w, "author name is required", http.StatusBadRequest)
		return
	}
	books := h.repo.List()
	for _, book := range books {
		for _, author := range book.Authors {
			if strings.TrimSpace(author.FullName) == name {
				http.Redirect(w, r, "/authors?error=published", http.StatusSeeOther)
				return
			}
		}
	}
	for _, book := range books {
		filtered := make([]model.Author, 0, len(book.Authors))
		for _, author := range book.Authors {
			if strings.TrimSpace(author.FullName) == name {
				continue
			}
			filtered = append(filtered, author)
		}
		if len(filtered) != len(book.Authors) {
			book.Authors = filtered
			h.repo.Update(book.ID, book)
		}
	}
	http.Redirect(w, r, "/authors", http.StatusSeeOther)
}

func authorBioValue(value string) *string {
	v := strings.TrimSpace(value)
	if v == "" {
		return nil
	}
	return &v
}

func collectAuthorCards(books []model.Book) []pages.AuthorCard {
	byAuthor := make(map[string]pages.AuthorCard)
	order := make([]string, 0)
	for _, book := range books {
		for _, author := range book.Authors {
			name := strings.TrimSpace(author.FullName)
			if name == "" {
				continue
			}
			if _, exists := byAuthor[name]; !exists {
				order = append(order, name)
				byAuthor[name] = pages.AuthorCard{Name: name, Books: nil}
			}
			card := byAuthor[name]
			card.Books = append(card.Books, pages.AuthorBook{ID: book.ID, Title: book.Title, URL: "/books/" + strconv.Itoa(book.ID) + "/edit"})
			if author.Bio != nil {
				if bio := strings.TrimSpace(*author.Bio); bio != "" {
					card.Bio = bio
				}
			}
			if strings.TrimSpace(card.Bio) == "" {
				card.Bio = authorBio(author)
			}
			byAuthor[name] = card
		}
	}
	cards := make([]pages.AuthorCard, 0, len(order))
	for _, name := range order {
		card := byAuthor[name]
		sort.Slice(card.Books, func(i, j int) bool { return card.Books[i].Title < card.Books[j].Title })
		cards = append(cards, card)
	}
	return cards
}

func parseAuthorNameFromPath(rawPath string, suffix string) string {
	trimmed := strings.Trim(rawPath, "/")
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "authors/") {
		trimmed = strings.TrimPrefix(trimmed, "authors/")
	}
	if suffix != "" && strings.HasSuffix(trimmed, "/"+suffix) {
		trimmed = strings.TrimSuffix(trimmed, "/"+suffix)
	}
	decoded, err := url.QueryUnescape(trimmed)
	if err != nil {
		decoded = trimmed
	}
	decoded = strings.TrimSpace(decoded)
	if decoded == "" {
		return ""
	}
	return decoded
}

func authorBio(author model.Author) string {
	if author.Bio != nil {
		bio := strings.TrimSpace(*author.Bio)
		if bio != "" {
			return bio
		}
	}
	return "This author is present in the catalog and starts a new reading path with each new title."
}

func matchesAuthorQuery(name string, books []pages.AuthorBook, bio string, query string) bool {
	needle := strings.ToLower(strings.TrimSpace(query))
	if needle == "" {
		return true
	}
	if strings.Contains(strings.ToLower(name), needle) {
		return true
	}
	for _, book := range books {
		if strings.Contains(strings.ToLower(book.Title), needle) {
			return true
		}
	}
	return strings.Contains(strings.ToLower(bio), needle)
}
