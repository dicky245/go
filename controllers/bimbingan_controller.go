package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"

	"time"
	"fmt"
)

// Ambil semua bimbingan milik user (berdasarkan token)
func GetBimbingan(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id tidak ditemukan di token"})
		return
	}

	userID := userIDInterface.(uint)

	// Cari data kelompok_id dari tabel kelompok_mahasiswa berdasarkan user_id
	var km model.KelompokMahasiswa
	if err := config.DB.Where("user_id = ?", userID).First(&km).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kelompok tidak ditemukan untuk user"})
		return
	}

	// Gunakan km.KelompokID untuk cari semua bimbingan
	var bimbingans []model.Bimbingan
	if err := config.DB.
		Where("kelompok_id = ?", km.KelompokID).
		Preload("Kelompok").
		Find(&bimbingans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data bimbingan"})
		return
	}
	

	c.JSON(http.StatusOK, bimbingans)
}

// Tambah request bimbingan (hanya mahasiswa)
func CreateBimbingan(c *gin.Context) {
	// Ambil user_id dari token
	userID := c.MustGet("user_id").(uint)

	// Ambil kelompok_id dari user_id
	var km model.KelompokMahasiswa
	if err := config.DB.Where("user_id = ?", userID).First(&km).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Mahasiswa belum tergabung dalam kelompok",
		})
		return
	}

	// Bind request body ke struct Bimbingan
	var req model.Bimbingan
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("BindJSON error:", err.Error()) // <--- tambahkan ini
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	

	// Set user_id dan kelompok_id pada permintaan bimbingan
	req.UserID = userID
	req.KelompokID = km.KelompokID
	req.Status = "menunggu" // default status
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	// Simpan bimbingan ke database
	if err := config.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create bimbingan",
			"details": err.Error(),
		})
		return
	}

	// Kembalikan response
	c.JSON(http.StatusCreated, req)
}
