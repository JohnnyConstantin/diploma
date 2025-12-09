package auth

import (
	"diploma/pkg/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Claims JWT-претензии, использующиеся для авторизации пользователя.
type Claims struct {
	jwt.RegisteredClaims
	Login string `json:"login"`
}

// LoginRequest тело запроса на регистрацию или вход пользователя.
type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SetTokenCookie создает JWT-токен в куки и хедер authorization
func SetTokenCookie(c *gin.Context, cfg *config.Config, w http.ResponseWriter, login string) (string, error) {
	token, err := GenerateJWTToken(login, cfg)
	if err != nil {
		return "", err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     cfg.JWTCookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(cfg.JWTTokenExp),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   cfg.Env == "prod",
	})
	c.Header("Authorization", "Bearer "+token)
	return token, nil
}
