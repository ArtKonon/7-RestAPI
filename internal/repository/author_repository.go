package repository

import (
	"context"

	"libratrack/internal/model"

	"gorm.io/gorm"
)

type AuthorRepository interface {
	Create(ctx context.Context, author *model.Author) error
	List(ctx context.Context) ([]model.Author, error)
	FindByID(ctx context.Context, id uint) (*model.Author, error)
}

type gormAuthorRepository struct {
	db *gorm.DB
}

func NewAuthorRepositoryWithDB(db *gorm.DB) AuthorRepository {
	return &gormAuthorRepository{db: db}
}

func (r *gormAuthorRepository) Create(ctx context.Context, author *model.Author) error {
	return r.db.WithContext(ctx).Create(author).Error
}

func (r *gormAuthorRepository) List(ctx context.Context) ([]model.Author, error) {
	var out []model.Author
	if err := r.db.WithContext(ctx).Order("id").Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *gormAuthorRepository) FindByID(ctx context.Context, id uint) (*model.Author, error) {
	var author model.Author
	if err := r.db.WithContext(ctx).First(&author, id).Error; err != nil {
		return nil, err
	}
	return &author, nil
}
