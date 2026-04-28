package user

type CreateUserRequest struct {
	UserNo   string `json:"user_no" binding:"required"`
	Name     string `json:"name" binding:"required"`
	TenantID string `json:"tenant_id" binding:"required"`
	Role     string `json:"role"`
	Password string `json:"password"`
}
