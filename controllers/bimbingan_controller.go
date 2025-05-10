package controllers

import (
	"net/http"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
)

// Ambil semua bimbingan milik user (berdasarkan token)
func GetBimbingan(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id tidak ditemukan di token"})
		return
	}

	userID := userIDInterface.(uint)

	var km model.KelompokMahasiswa
	if err := config.DB.Where("user_id = ?", userID).First(&km).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Mahasiswa belum memiliki kelompok",
			"status":  "no_group",
			"data":    []map[string]interface{}{},
		})
		return
	}

	var bimbingans []model.Bimbingan
	if err := config.DB.
		Where("kelompok_id = ?", km.KelompokID).
		Preload("Kelompok").
		Preload("Ruangan"). // preload ruangan juga
		Find(&bimbingans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data bimbingan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data bimbingan berhasil diambil",
		"status":  "success",
		"data":    bimbingans,
	})
}

// Tambah request bimbingan (hanya mahasiswa)
func CreateBimbingan(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var km model.KelompokMahasiswa
	if err := config.DB.Where("user_id = ?", userID).First(&km).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Mahasiswa belum tergabung dalam kelompok",
			"status":  "no_group",
		})
		return
	}

	var req model.Bimbingan
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("BindJSON error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userID
	req.KelompokID = km.KelompokID
	req.Status = "menunggu"
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	if err := config.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create bimbingan",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bimbingan berhasil dibuat",
		"status":  "success",
		"data":    req,
	})
}
