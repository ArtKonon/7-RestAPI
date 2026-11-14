package repository

import (
	"context"
	"strings"

	"libratrack/internal/model"

	"gorm.io/gorm"
)

type BookRepository interface {
	Create(book model.Book) model.Book
	GetByID(id int) (model.Book, bool)
	Update(id int, newBook model.Book) (model.Book, bool)
	Delete(id int) bool
	List() []model.Book
}

type gormBookRepository struct {
	db *gorm.DB
}

func NewGormBookRepository(db *gorm.DB) BookRepository {
	return &gormBookRepository{db: db}
}

func (r *gormBookRepository) Create(book model.Book) model.Book {
	if r.db == nil {
		return book
	}
	ctx := context.Background()
	if category, err := r.ensureCategory(ctx, strings.TrimSpace(book.Category.Name)); err == nil && category.ID > 0 {
		book.CategoryID = category.ID
		book.Category = category
	}
	if authors, err := r.ensureAuthors(ctx, book.Authors); err == nil && len(authors) > 0 {
		book.Authors = authors
	}
	if err := r.db.WithContext(ctx).Create(&book).Error; err != nil {
		return book
	}
	if len(book.Authors) > 0 {
		_ = r.db.WithContext(ctx).Model(&book).Association("Authors").Replace(book.Authors)
	}
	if book.CategoryID > 0 {
		_ = r.db.WithContext(ctx).Model(&book).Update("category_id", book.CategoryID)
	}
	_ = r.db.WithContext(ctx).Preload("Authors").Preload("Category").First(&book, book.ID).Error
	return book
}

func (r *gormBookRepository) GetByID(id int) (model.Book, bool) {
	var book model.Book
	if r.db == nil {
		return model.Book{}, false
	}
	if err := r.db.WithContext(context.Background()).Preload("Authors").Preload("Category").First(&book, id).Error; err != nil {
		return model.Book{}, false
	}
	return book, true
}

func (r *gormBookRepository) Update(id int, newBook model.Book) (model.Book, bool) {
	if r.db == nil {
		return model.Book{}, false
	}
	ctx := context.Background()
	var existing model.Book
	if err := r.db.WithContext(ctx).Preload("Authors").Preload("Category").First(&existing, id).Error; err != nil {
		return model.Book{}, false
	}

	if strings.TrimSpace(newBook.Title) != "" {
		existing.Title = newBook.Title
	}
	if strings.TrimSpace(newBook.ISBN) != "" {
		existing.ISBN = newBook.ISBN
	}
	if newBook.Description != nil {
		existing.Description = newBook.Description
	}
	if strings.TrimSpace(newBook.Category.Name) != "" {
		category, err := r.ensureCategory(ctx, newBook.Category.Name)
		if err == nil && category.ID > 0 {
			existing.CategoryID = category.ID
			existing.Category = category
		}
	}
	if len(newBook.Authors) > 0 {
		authors, err := r.ensureAuthors(ctx, newBook.Authors)
		if err == nil {
			existing.Authors = authors
		}
	}

	if err := r.db.WithContext(ctx).Model(&existing).Updates(map[string]interface{}{
		"title":       existing.Title,
		"isbn":        existing.ISBN,
		"category_id": existing.CategoryID,
		"description": existing.Description,
	}).Error; err != nil {
		return model.Book{}, false
	}
	if len(existing.Authors) > 0 {
		if err := r.db.WithContext(ctx).Model(&existing).Association("Authors").Replace(existing.Authors); err != nil {
			return model.Book{}, false
		}
	}
	if err := r.db.WithContext(ctx).Preload("Authors").Preload("Category").First(&existing, id).Error; err != nil {
		return model.Book{}, false
	}
	return existing, true
}

func (r *gormBookRepository) ensureCategory(ctx context.Context, name string) (model.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return model.Category{}, nil
	}
	var category model.Category
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&category).Error; err == nil {
		return category, nil
	}
	category.Name = name
	if err := r.db.WithContext(ctx).Create(&category).Error; err != nil {
		return model.Category{}, err
	}
	return category, nil
}

func normalizeAuthorBio(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func (r *gormBookRepository) ensureAuthors(ctx context.Context, authors []model.Author) ([]model.Author, error) {
	if len(authors) == 0 {
		return nil, nil
	}
	result := make([]model.Author, 0, len(authors))
	for _, author := range authors {
		name := strings.TrimSpace(author.FullName)
		if name == "" {
			continue
		}

		incomingBio := normalizeAuthorBio(author.Bio)
		var existing model.Author
		if err := r.db.WithContext(ctx).Where("full_name = ?", name).First(&existing).Error; err == nil {
			updatedBio := normalizeAuthorBio(existing.Bio)
			if incomingBio != nil {
				updatedBio = incomingBio
			} else if existing.Bio != nil {
				updatedBio = nil
			}
			if existing.FullName != name || !authorBioEqual(existing.Bio, updatedBio) {
				existing.FullName = name
				existing.Bio = updatedBio
				if err := r.db.WithContext(ctx).Model(&existing).Updates(map[string]interface{}{
					"full_name": existing.FullName,
					"bio":       existing.Bio,
				}).Error; err != nil {
					return nil, err
				}
			}
			result = append(result, existing)
			continue
		}

		existing.FullName = name
		existing.Bio = incomingBio
		if err := r.db.WithContext(ctx).Create(&existing).Error; err != nil {
			return nil, err
		}
		result = append(result, existing)
	}
	return result, nil
}

func authorBioEqual(current *string, next *string) bool {
	if current == nil && next == nil {
		return true
	}
	if current == nil || next == nil {
		return false
	}
	return strings.TrimSpace(*current) == strings.TrimSpace(*next)
}

func (r *gormBookRepository) Delete(id int) bool {
	if r.db == nil {
		return false
	}
	return r.db.WithContext(context.Background()).Delete(&model.Book{}, id).Error == nil
}

func (r *gormBookRepository) List() []model.Book {
	var books []model.Book
	if r.db == nil {
		return nil
	}
	if err := r.db.WithContext(context.Background()).Preload("Authors").Preload("Category").Order("id").Find(&books).Error; err != nil {
		return nil
	}
	return books
}
