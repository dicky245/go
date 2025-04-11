package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
)

func GetPengumuman(c *gin.Context) {
	db := config.DB

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	var pengumumanList []model.Pengumuman
	if err := db.Where("user_id = ?", userID.(uint)).Find(&pengumumanList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pengumumanList)
}

// Ambil satu pengumuman berdasarkan ID (dan user yang login)
func GetPengumumanByID(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	var pengumuman model.Pengumuman
	if err := db.Where("pengumuman_id = ? AND user_id = ?", id, userID).First(&pengumuman).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengumuman tidak ditemukan atau bukan milik Anda"})
		return
	}

	c.JSON(http.StatusOK, pengumuman)
}

// Tambah request pengumuman (hanya mahasiswa)
func CreatePengumuman(c *gin.Context) {
	db := config.DB
	var newPengumuman model.Pengumuman

	// Ambil user_id dan role dari token
	userID, exists := c.Get("user_id")
	role, roleExists := c.Get("user_role")
	if !exists || !roleExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	// Hanya mahasiswa yang bisa membuat Pengumuman
	if role != "Mahasiswa" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Mahasiswa yang dapat mengajukan Pengumuman"})
		return
	}

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&newPengumuman); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set user_id dari token
	newPengumuman.UserID = userID.(uint)

	// Simpan ke DB
	if err := db.Create(&newPengumuman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pengumuman berhasil ditambahkan", "data": newPengumuman})
}
func UpdatePengumuman(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	var pengumuman model.Pengumuman
	if err := db.Where("pengumuman_id = ?", id).First(&pengumuman).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengumuman tidak ditemukan"})
		return
	}

	// Pastikan user yang login adalah pemiliknya
	userID, exists := c.Get("user_id")
	if !exists || pengumuman.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk mengubah pengumuman ini"})
		return
	}

	// Bind JSON
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update data
	if err := db.Model(&pengumuman).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengumuman berhasil diperbarui", "data": pengumuman})
}

// Hapus request pengumuman
func DeletePengumuman(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	var pengumuman model.Pengumuman
	if err := db.Where("pengumuman_id = ?", id).First(&pengumuman).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengumuman tidak ditemukan"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists || pengumuman.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk menghapus pengumuman ini"})
		return
	}

	if err := db.Delete(&pengumuman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengumuman berhasil dihapus"})
}
