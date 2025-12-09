package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// заглушка
func newTestAPI() *API {
	return &API{}
}

// Общая функция для создания контекста для всех тестов (чтобы не повторять имплементацию для каждого отдельного случая)
func newTestContext(method, target string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(method, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

func TestPostPassword_UnauthorizedWhenNoLoginInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	api := newTestAPI()

	c, w := newTestContext(http.MethodPost, "/api/v1/passwords", []byte(`{}`))

	api.PostPassword(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestPostPassword_BadRequestOnInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	api := newTestAPI()

	c, w := newTestContext(http.MethodPost, "/api/v1/passwords", []byte(`{invalid json`))

	c.Set("user_login", "testuser")

	api.PostPassword(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetPasswordHandler_UnauthorizedWhenNoLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	api := newTestAPI()

	c, w := newTestContext(http.MethodGet, "/api/v1/passwords/some-id", nil)

	api.GetPasswordHandler(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}
