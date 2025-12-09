package auth

import (
	"diploma/pkg/config"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	// ErrNoAuthHeader возвращается, если заголовок Authorization отсутствует.
	ErrNoAuthHeader = errors.New("no auth header")
	// ErrBadAuthHeader возвращается, если заголовок Authorization имеет некорректный формат.
	ErrBadAuthHeader = errors.New("bad auth header")
	// ErrNoCookie возвращается, если JWT-токен не найден в cookie.
	ErrNoCookie = errors.New("no cookie")
	// ErrUnexpectedSigningMethod возвращается, если токен подписан неожиданным методом.
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	// ErrInvalidToken возвращается, если JWT-токен невалиден или подпись не сходится.
	ErrInvalidToken = errors.New("invalid token")
)

// GenerateJWTToken создаёт и подписывает JWT-токен для указанного логина
// с параметрами из конфигурации.
func GenerateJWTToken(login string, cfg *config.Config) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWTTokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Login: login,
	})
	return token.SignedString([]byte(cfg.JWTSecretKey))
}

func parseTokenFromCookie(r *http.Request, cfg *config.Config) (string, error) {
	c, err := r.Cookie(cfg.JWTCookieName)
	if err != nil || c == nil || c.Value == "" {
		return "", ErrNoCookie
	}
	return c.Value, nil
}

func parseTokenFromAuthHeader(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return "", ErrNoAuthHeader
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", ErrBadAuthHeader
	}
	return parts[1], nil
}

func parseLoginFromToken(raw string, cfg *config.Config) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(cfg.JWTSecretKey), nil
	})
	if err != nil || !token.Valid {
		return "", ErrInvalidToken
	}
	return claims.Login, nil
}
