package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
)

// Ambil semua bimbingan milik user (berdasarkan token)
func GetBimbingan(c *gin.Context) {
	db := config.DB

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	var bimbinganList []model.Bimbingan
	if err := db.Where("user_id = ?", userID.(uint)).Find(&bimbinganList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bimbinganList)
}

// Ambil satu bimbingan berdasarkan ID (dan user yang login)
func GetBimbinganByID(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	var bimbingan model.Bimbingan
	if err := db.Where("bimbingan_id = ? AND user_id = ?", id, userID).First(&bimbingan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bimbingan tidak ditemukan atau bukan milik Anda"})
		return
	}

	c.JSON(http.StatusOK, bimbingan)
}

// Tambah request bimbingan (hanya mahasiswa)
func CreateBimbingan(c *gin.Context) {
	db := config.DB
	var newBimbingan model.Bimbingan

	// Ambil user_id dan role dari token
	userID, exists := c.Get("user_id")
	role, roleExists := c.Get("user_role")
	if !exists || !roleExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	// Hanya mahasiswa yang bisa membuat bimbingan
	if role != "Mahasiswa" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Mahasiswa yang dapat mengajukan bimbingan"})
		return
	}

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&newBimbingan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set user_id dari token
	newBimbingan.UserID = userID.(uint)

	// Simpan ke DB
	if err := db.Create(&newBimbingan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Bimbingan berhasil ditambahkan", "data": newBimbingan})
}

// Update request bimbingan (hanya milik sendiri)
func UpdateBimbingan(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	var bimbingan model.Bimbingan
	if err := db.Where("bimbingan_id = ?", id).First(&bimbingan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bimbingan tidak ditemukan"})
		return
	}

	// Pastikan user yang login adalah pemiliknya
	userID, exists := c.Get("user_id")
	if !exists || bimbingan.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk mengubah bimbingan ini"})
		return
	}

	// Bind JSON
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update data
	if err := db.Model(&bimbingan).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bimbingan berhasil diperbarui", "data": bimbingan})
}

// Hapus request bimbingan
func DeleteBimbingan(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	var bimbingan model.Bimbingan
	if err := db.Where("bimbingan_id = ?", id).First(&bimbingan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bimbingan tidak ditemukan"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists || bimbingan.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk menghapus bimbingan ini"})
		return
	}

	if err := db.Delete(&bimbingan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bimbingan berhasil dihapus"})
}
