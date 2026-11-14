package web

import (
	"net/http"
	"sort"
	"strings"

	"libratrack/internal/model"
	pages "libratrack/web/templates/pages"
)

func (h *Handler) AuthorsPage(w http.ResponseWriter, r *http.Request) {
	profile := readProfileCookie(r)
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	books := h.repo.List()
	byAuthor := make(map[string]pages.AuthorCard)
	order := make([]string, 0)
	for _, book := range books {
		if len(book.Authors) == 0 {
			continue
		}
		for _, author := range book.Authors {
			name := strings.TrimSpace(author.FullName)
			if name == "" {
				continue
			}
			if _, exists := byAuthor[name]; !exists {
				order = append(order, name)
				byAuthor[name] = pages.AuthorCard{Name: name, Bio: authorBio(author), Books: nil}
			}
			card := byAuthor[name]
			card.Books = append(card.Books, book.Title)
			if strings.TrimSpace(card.Bio) == "" {
				card.Bio = authorBio(author)
			}
			byAuthor[name] = card
		}
	}
	cards := make([]pages.AuthorCard, 0, len(order))
	for _, name := range order {
		card := byAuthor[name]
		if query != "" && !matchesAuthorQuery(name, card.Books, card.Bio, query) {
			continue
		}
		sort.Strings(card.Books)
		cards = append(cards, card)
	}
	page := pages.AuthorsPage(cards, profile.Name, query)
	h.renderPage(w, r, page, pages.AuthorsContent(cards, profile.Name, query), http.StatusOK)
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

func matchesAuthorQuery(name string, books []string, bio string, query string) bool {
	needle := strings.ToLower(strings.TrimSpace(query))
	if needle == "" {
		return true
	}
	if strings.Contains(strings.ToLower(name), needle) {
		return true
	}
	for _, title := range books {
		if strings.Contains(strings.ToLower(title), needle) {
			return true
		}
	}
	return strings.Contains(strings.ToLower(bio), needle)
}
