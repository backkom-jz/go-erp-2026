package auth

type LoginRequest struct {
	UserNo   string `json:"user_no" binding:"required"`
	TenantID string `json:"tenant_id" binding:"required"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
