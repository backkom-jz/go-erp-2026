package user

import (
	"context"
	"go-erp/internal/domain/user"
	dtouser "go-erp/internal/dto/user"
	userrepo "go-erp/internal/repository/user"
	"go-erp/pkg/errs"
)

type Service struct {
	repo userrepo.Repository
}

const defaultInitPassword = "dev_init_password"

// NewService 创建用户服务。
func NewService(repo userrepo.Repository) *Service {
	return &Service{repo: repo}
}

// Create 创建用户。
func (s *Service) Create(ctx context.Context, req dtouser.CreateUserRequest) error {
	entity := &user.User{
		UserNo:   req.UserNo,
		Name:     req.Name,
		Password: req.Password,
		TenantID: req.TenantID,
		Role:     req.Role,
	}
	if entity.Role == "" {
		entity.Role = "viewer"
	}
	if entity.Password == "" {
		entity.Password = defaultInitPassword
	}
	return s.repo.Create(ctx, entity)
}

// GetByUserNo 通过 user_no 查询用户。
func (s *Service) GetByUserNo(ctx context.Context, userNo string) (*user.User, error) {
	entity, err := s.repo.GetByUserNo(ctx, userNo)
	if err != nil {
		return nil, errs.Wrap(errs.CodeNotFound, "user_not_found", err)
	}
	return entity, nil
}
