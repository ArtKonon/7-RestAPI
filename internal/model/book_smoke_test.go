package model

import (
	"encoding/json"
	"testing"
)

func TestBookModelSmoke(t *testing.T) {
	book := Book{
		Title:      "The Go Programming Language",
		ISBN:       "9780134190440",
		CategoryID: 7,
		Category:   Category{ID: 7, Name: "programming"},
		Authors:    []Author{{ID: 1, FullName: "Brian Kernighan"}},
	}

	if got := book.LegacyAuthor(); got != "Brian Kernighan" {
		t.Fatalf("LegacyAuthor() = %q; want %q", got, "Brian Kernighan")
	}
	if got := book.LegacyCategory(); got != "programming" {
		t.Fatalf("LegacyCategory() = %q; want %q", got, "programming")
	}

	payload, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("MarshalJSON() error: %v", err)
	}
	if len(payload) == 0 {
		t.Fatal("MarshalJSON() returned empty payload")
	}

	var decoded Book
	if err := json.Unmarshal([]byte(`{"title":"Test","isbn":"1","category_id":7,"author":"Alice","category":"programming"}`), &decoded); err != nil {
		t.Fatalf("UnmarshalJSON() error: %v", err)
	}
	if decoded.Category.Name != "programming" {
		t.Fatalf("Category name mismatch: got %q", decoded.Category.Name)
	}
	if len(decoded.Authors) != 1 || decoded.Authors[0].FullName != "Alice" {
		t.Fatalf("Author data was not normalized correctly: %#v", decoded.Authors)
	}
}
