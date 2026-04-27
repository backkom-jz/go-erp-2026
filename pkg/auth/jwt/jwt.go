package jwt

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
	jwtlib.RegisteredClaims
}

type Manager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewManager(secret string, accessTTLMinutes int, refreshTTLMinutes int) *Manager {
	if accessTTLMinutes <= 0 {
		accessTTLMinutes = 30
	}
	if refreshTTLMinutes <= 0 {
		refreshTTLMinutes = 60 * 24 * 7
	}
	return &Manager{
		secret:     []byte(secret),
		accessTTL:  time.Duration(accessTTLMinutes) * time.Minute,
		refreshTTL: time.Duration(refreshTTLMinutes) * time.Minute,
	}
}

func (m *Manager) SignAccessToken(userID, tenantID, role string) (string, error) {
	return m.signToken(userID, tenantID, role, m.accessTTL)
}

func (m *Manager) SignRefreshToken(userID, tenantID, role string) (string, error) {
	return m.signToken(userID, tenantID, role, m.refreshTTL)
}

func (m *Manager) signToken(userID, tenantID, role string, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Manager) Parse(token string) (*Claims, error) {
	parsed, err := jwtlib.ParseWithClaims(token, &Claims{}, func(t *jwtlib.Token) (interface{}, error) {
		if t.Method != jwtlib.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
