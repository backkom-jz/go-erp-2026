package auth

import (
	"context"
	dtouser "go-erp/internal/dto/user"
	domainuser "go-erp/internal/domain/user"
	userrepo "go-erp/internal/repository/user"
	"go-erp/pkg/auth/jwt"
	"go-erp/pkg/errs"
)

type Service struct {
	users userrepo.Repository
	jwt   *jwt.Manager
}

const defaultInitPassword = "dev_init_password"

// NewService 创建认证服务。
func NewService(users userrepo.Repository, jwtManager *jwt.Manager) *Service {
	return &Service{
		users: users,
		jwt:   jwtManager,
	}
}

// Login 处理登录逻辑。
// 备注：若用户不存在会自动初始化一个默认用户。
func (s *Service) Login(ctx context.Context, userNo, tenantID, role string) (string, string, error) {
	u, err := s.users.GetByUserNo(ctx, userNo)
	if err != nil {
		if role == "" {
			role = "viewer"
		}
		newUser := &domainuser.User{
			UserNo:   userNo,
			Name:     userNo,
			Password: defaultInitPassword,
			TenantID: tenantID,
			Role:     role,
		}
		if createErr := s.users.Create(ctx, newUser); createErr != nil {
			return "", "", errs.Wrap(errs.CodeInternal, "create_user_failed", createErr)
		}
		u = newUser
	}

	accessToken, err := s.jwt.SignAccessToken(u.UserNo, tenantID, u.Role)
	if err != nil {
		return "", "", errs.Wrap(errs.CodeInternal, "sign_access_token_failed", err)
	}
	refreshToken, err := s.jwt.SignRefreshToken(u.UserNo, tenantID, u.Role)
	if err != nil {
		return "", "", errs.Wrap(errs.CodeInternal, "sign_refresh_token_failed", err)
	}
	return accessToken, refreshToken, nil
}

// Register 注册用户（用于管理端调用）。
func (s *Service) Register(ctx context.Context, req dtouser.CreateUserRequest) error {
	entity := &domainuser.User{
		UserNo:   req.UserNo,
		Name:     req.Name,
		Password: req.Password,
		TenantID: req.TenantID,
		Role:     req.Role,
	}
	if entity.Password == "" {
		entity.Password = defaultInitPassword
	}
	return s.users.Create(ctx, entity)
}
