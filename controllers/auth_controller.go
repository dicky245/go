package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/rudychandra/lagi/utils"

	"github.com/gin-gonic/gin"
)

// Struktur untuk request login
type LoginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// Struktur response dari API eksternal (CIS)
type LoginResponse struct {
	Result bool   `json:"result"`
	Error  string `json:"error,omitempty"`
	Token  string `json:"token,omitempty"`
	User   struct {
		UserID   int    `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	} `json:"user,omitempty"`
}

// Fungsi login ke API eksternal dan generate token internal
func Login(c *gin.Context) {
	var loginReq LoginRequest

	// Bind request dari form-data
	if err := c.ShouldBind(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format form-data salah"})
		return
	}

	// URL API eksternal CIS
	apiURL := "https://cis-dev.del.ac.id/api/jwt-api/do-auth"

	// Buat body form-data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	_ = writer.WriteField("username", loginReq.Username)
	_ = writer.WriteField("password", loginReq.Password)
	writer.Close()

	// Buat HTTP request
	req, err := http.NewRequest("POST", apiURL, &requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat request ke server autentikasi"})
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Kirim request ke API eksternal
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal terhubung ke server autentikasi"})
		return
	}
	defer resp.Body.Close()

	// Baca body response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca respons dari server autentikasi"})
		return
	}

	// Debug
	fmt.Println("[DEBUG] Response Body:", string(body))

	// Decode response
	var loginRes LoginResponse
	if err := json.Unmarshal(body, &loginRes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Format respons tidak valid dari server autentikasi"})
		return
	}

	// Jika login gagal
	if !loginRes.Result {
		c.JSON(http.StatusUnauthorized, gin.H{"error": loginRes.Error})
		return
	}

	// Generate token internal
	internalToken, err := utils.GenerateInternalToken(uint(loginRes.User.UserID), loginRes.User.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token internal"})
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, gin.H{
		"message":        "Login berhasil",
		"external_token": loginRes.Token,
		"internal_token": internalToken,
		"user": gin.H{
			"user_id":  loginRes.User.UserID,
			"username": loginRes.User.Username,
			"email":    loginRes.User.Email,
			"role":     loginRes.User.Role,
		},
	})
}
