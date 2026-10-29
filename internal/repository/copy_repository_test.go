package repository

import (
	"context"
	"testing"

	"libratrack/internal/model"
)

func TestCopyRepositorySmoke(t *testing.T) {
	repo := NewCopyRepositoryWithDB(nil)
	if repo == nil {
		t.Fatal("NewCopyRepositoryWithDB returned nil")
	}

	if err := repo.Create(context.Background(), &model.Copy{InventoryNo: "INV-001", Status: "available"}); err != nil {
		t.Fatalf("Create() returned unexpected error: %v", err)
	}

	items, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("List() returned unexpected error: %v", err)
	}
	if items != nil {
		t.Fatalf("List() with nil DB should return nil slice, got %#v", items)
	}
}
