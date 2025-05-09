package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// StudentResponse represents the structure of the student data from external API
type StudentResponse struct {
	Result string `json:"result"`
	Data   struct {
		Mahasiswa []struct {
			DimID     int    `json:"dim_id"`
			UserID    int    `json:"user_id"`
			UserName  string `json:"user_name"`
			Nim       string `json:"nim"`
			Nama      string `json:"nama"`
			Email     string `json:"email"`
			ProdiID   int    `json:"prodi_id"`
			ProdiName string `json:"prodi_name"`
			Fakultas  string `json:"fakultas"`
			Angkatan  int    `json:"angkatan"`
			Status    string `json:"status"`
			Asrama    string `json:"asrama"`
		} `json:"mahasiswa"`
	} `json:"data"`
}

// GetStudentData fetches student data from external API
func GetStudentData(c *gin.Context) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
		return
	}

	// Extract token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid"})
		return
	}
	token := tokenParts[1]

	// Get username from query parameter
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter username diperlukan"})
		return
	}

	// Construct URL for external API
	externalAPIURL := "https://cis-dev.del.ac.id/api/library-api/mahasiswa"
	req, err := http.NewRequest("GET", externalAPIURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat request ke API eksternal"})
		return
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("username", username)
	req.URL.RawQuery = q.Encode()

	// Add authorization header
	req.Header.Set("Authorization", "Bearer "+token)

	// Send request to external API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal terhubung ke API eksternal"})
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "API eksternal mengembalikan status error"})
		return
	}

	// Parse response
	var studentResp StudentResponse
	if err := json.NewDecoder(resp.Body).Decode(&studentResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses respons dari API eksternal"})
		return
	}

	// Check if student data exists
	if studentResp.Result != "Ok" || len(studentResp.Data.Mahasiswa) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data mahasiswa tidak ditemukan"})
		return
	}

	// Return student data
	c.JSON(http.StatusOK, studentResp)
}