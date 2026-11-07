package middleware

import (
	"net/http"
	"strings"

	"storage-manager-api/internal/utils"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey = "user_id"
	ContextEmailKey  = "email"
)

// AuthMiddleware memvalidasi header Authorization Bearer JWT dan menyimpan klaim user ke context Gin.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Header Authorization tidak ditemukan")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Format Authorization harus 'Bearer <token>'")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(strings.TrimSpace(parts[1]), jwtSecret)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Token tidak valid atau sudah kedaluwarsa")
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextEmailKey, claims.Email)
		c.Next()
	}
}
