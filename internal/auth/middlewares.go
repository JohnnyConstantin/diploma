package auth

import (
	"context"
	"diploma/internal/repo"
	"diploma/pkg/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func isPublicPath(path string) bool {
	public := []string{
		"/health",
		"/api/user/register",
		"/api/user/login",
	}
	if strings.HasPrefix(path, "/api/secret/") {
		return true
	}
	for _, p := range public {
		if path == p {
			return true
		}
	}
	return path == "health"
}

// MiddlewareAuth проверяет JWT-токен в cookie или заголовке Authorization,
// извлекает логин пользователя и кладёт его в контекст запроса.
func MiddlewareAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isPublicPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		var rawToken string
		var login string
		var err error

		if t, e := parseTokenFromCookie(c.Request, cfg); e == nil {
			rawToken = t
		} else {
			if t, e2 := parseTokenFromAuthHeader(c.Request); e2 == nil {
				rawToken = t
			}
		}

		if rawToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		login, err = parseLoginFromToken(rawToken, cfg)
		if err != nil || strings.TrimSpace(login) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Скользящее продление.
		_, _ = SetTokenCookie(c, cfg, c.Writer, login)

		ctx := context.WithValue(c.Request.Context(), repo.UserLoginKey, login)
		c.Request = c.Request.WithContext(ctx)
		c.Set(string(repo.UserLoginKey), login)

		c.Next()
	}
}
