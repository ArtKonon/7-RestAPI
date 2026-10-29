package repository

import (
	"context"

	"libratrack/internal/model"

	"gorm.io/gorm"
)

type CopyRepository interface {
	Create(ctx context.Context, copy *model.Copy) error
	List(ctx context.Context) ([]model.Copy, error)
	FindByID(ctx context.Context, id uint) (*model.Copy, error)
	Update(ctx context.Context, copy *model.Copy) error
	Delete(ctx context.Context, id uint) error
}

type gormCopyRepository struct {
	db *gorm.DB
}

func NewCopyRepositoryWithDB(db *gorm.DB) CopyRepository {
	return &gormCopyRepository{db: db}
}

func (r *gormCopyRepository) Create(ctx context.Context, copy *model.Copy) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Create(copy).Error
}

func (r *gormCopyRepository) List(ctx context.Context) ([]model.Copy, error) {
	var out []model.Copy
	if r.db == nil {
		return out, nil
	}
	if err := r.db.WithContext(ctx).Order("id").Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *gormCopyRepository) FindByID(ctx context.Context, id uint) (*model.Copy, error) {
	var copy model.Copy
	if r.db == nil {
		return nil, gorm.ErrRecordNotFound
	}
	if err := r.db.WithContext(ctx).First(&copy, id).Error; err != nil {
		return nil, err
	}
	return &copy, nil
}

func (r *gormCopyRepository) Update(ctx context.Context, copy *model.Copy) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Save(copy).Error
}

func (r *gormCopyRepository) Delete(ctx context.Context, id uint) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Delete(&model.Copy{}, id).Error
}
