package jwt

import "testing"

func TestManagerSignAndParse(t *testing.T) {
	m := NewManager("secret", 10, 20)
	token, err := m.SignAccessToken("u100", "t100", "admin")
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}
	claims, err := m.Parse(token)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if claims.UserID != "u100" || claims.TenantID != "t100" || claims.Role != "admin" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}
