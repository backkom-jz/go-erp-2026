package testutil

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func NewGinEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func NewRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func MustEqual(t *testing.T, got interface{}, want interface{}) {
	t.Helper()
	if got != want {
		t.Fatalf("got=%v want=%v", got, want)
	}
}
