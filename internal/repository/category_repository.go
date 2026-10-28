package repository

import (
	"context"

	"libratrack/internal/model"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	List(ctx context.Context) ([]model.Category, error)
	FindByID(ctx context.Context, id uint) (*model.Category, error)
}

type gormCategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepositoryWithDB(db *gorm.DB) CategoryRepository {
	return &gormCategoryRepository{db: db}
}

func (r *gormCategoryRepository) Create(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *gormCategoryRepository) List(ctx context.Context) ([]model.Category, error) {
	var out []model.Category
	if err := r.db.WithContext(ctx).Order("id").Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *gormCategoryRepository) FindByID(ctx context.Context, id uint) (*model.Category, error) {
	var category model.Category
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
