package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-erp/pkg/auth/jwt"

	"github.com/gin-gonic/gin"
)

func TestJWTMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	m := jwt.NewManager("test_secret", 10, 20)
	token, err := m.SignAccessToken("u1", "t1", "admin")
	if err != nil {
		t.Fatalf("sign token failed: %v", err)
	}

	r := gin.New()
	r.Use(JWT(m))
	r.GET("/api/v1/users/me", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
