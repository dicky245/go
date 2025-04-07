package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/utils"
)

func InternalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token internal tidak valid"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Verifikasi token menggunakan fungsi dari utils
		claims, err := utils.VerifyInternalToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah kadaluarsa"})
			c.Abort()
			return
		}

		// Simpan data user ke context
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}
