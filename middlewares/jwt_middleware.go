package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Struktur response dari API eksternal
type AuthResponse struct {
	Result bool   `json:"result"`
	Error  string `json:"error,omitempty"`
	User   struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	} `json:"user,omitempty"`
}

// Middleware untuk autentikasi dengan API eksternal
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			c.Abort()
			return
		}

		// Pastikan format token valid
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid"})
			c.Abort()
			return
		}
		token := tokenParts[1]

		// Kirim token ke API eksternal untuk verifikasi
		apiURL := "https://cis-dev.del.ac.id/api/jwt-api/verify-token"
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat request ke server autentikasi"})
			c.Abort()
			return
		}
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah kadaluarsa"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		// Parse response dari API eksternal
		var authRes AuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&authRes); err != nil || !authRes.Result {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Gagal memverifikasi token"})
			c.Abort()
			return
		}

		// Simpan user ID dalam context request
		c.Set("user_id", authRes.User.UserID)
		c.Set("user_role", authRes.User.Role)
		c.Next()
	}
}
