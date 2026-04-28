package user

import (
	"context"
	"go-erp/internal/domain/user"

	"gorm.io/gorm"
)

type Repository interface {
	// Create 创建用户记录。
	Create(ctx context.Context, entity *user.User) error
	// GetByUserNo 按 user_no 查询用户。
	GetByUserNo(ctx context.Context, userNo string) (*user.User, error)
	// GetByID 按主键 ID 查询用户。
	GetByID(ctx context.Context, id uint) (*user.User, error)
}

type GormRepository struct {
	db *gorm.DB
}

// NewRepository 创建用户仓储实现。
func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// Create 创建用户记录。
func (r *GormRepository) Create(ctx context.Context, entity *user.User) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// GetByUserNo 按 user_no 查询用户。
func (r *GormRepository) GetByUserNo(ctx context.Context, userNo string) (*user.User, error) {
	var entity user.User
	if err := r.db.WithContext(ctx).Where("user_no = ?", userNo).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// GetByID 按主键 ID 查询用户。
func (r *GormRepository) GetByID(ctx context.Context, id uint) (*user.User, error) {
	var entity user.User
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}
