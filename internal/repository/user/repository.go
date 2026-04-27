package user

import (
	"context"
	"go-erp/internal/domain/user"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, entity *user.User) error
	GetByUserNo(ctx context.Context, userNo string) (*user.User, error)
	GetByID(ctx context.Context, id uint) (*user.User, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, entity *user.User) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *GormRepository) GetByUserNo(ctx context.Context, userNo string) (*user.User, error) {
	var entity user.User
	if err := r.db.WithContext(ctx).Where("user_no = ?", userNo).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uint) (*user.User, error) {
	var entity user.User
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}
