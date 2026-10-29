package handler

import (
	"testing"
)

func TestBookHandlerSmoke(t *testing.T) {
	h := NewBookHandler(nil)
	if h == nil {
		t.Fatal("NewBookHandler returned nil")
	}

	id, err := parseID("42")
	if err != nil {
		t.Fatalf("parseID returned error: %v", err)
	}
	if id != 42 {
		t.Fatalf("parseID(42) = %d; want 42", id)
	}
}
