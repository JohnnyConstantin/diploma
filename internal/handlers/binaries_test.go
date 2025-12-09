package handlers

import (
	"diploma/internal/models"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestBinaryAPI() *API {
	return &API{}
}

func TestPostBinary_UnauthorizedWhenNoLoginInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	api := newTestBinaryAPI()

	c, w := newTestContext(http.MethodPost, "/api/v1/binaries", []byte(`{}`))

	api.PostBinary(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestPostBinary_BadRequestOnInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	api := newTestBinaryAPI()

	c, w := newTestContext(http.MethodPost, "/api/v1/binaries", []byte(`{invalid json`))

	c.Set("user_login", "testuser")

	api.PostBinary(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestPostBinary_TooLargeRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	api := newTestBinaryAPI()

	const MaxBytes = 101 * 1024 * 1024

	// Создаем запрос с Data > MaxBytes
	largeData := strings.Repeat("x", MaxBytes)

	req := models.BinaryRequest{
		Data: largeData,
		ID:   "2",
	}

	// Маршалим в JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	c, w := newTestContext(http.MethodPost, "/api/v1/binaries", jsonData)
	c.Set("user_login", "testuser")

	api.PostBinary(c)
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestGetBinaryHandler_UnauthorizedWhenNoLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	api := newTestBinaryAPI()

	c, w := newTestContext(http.MethodGet, "/api/v1/binaries/some-id", nil)

	api.GetBinaryHandler(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}
