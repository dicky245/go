package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
)

func GetArtefak(c *gin.Context) {
	db := config.DB

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	var artefakList []model.Artefak
	if err := db.Where("user_id = ?", userID.(uint)).Find(&artefakList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, artefakList)
}

func GetArtefakByID(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	var artefak model.Artefak
	if err := db.Where("artefak_id = ? AND user_id = ?", id, userID).First(&artefak).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artefak tidak ditemukan atau bukan milik Anda"})
		return
	}

	c.JSON(http.StatusOK, artefak)
}

func CreateArtefak(c *gin.Context) {
	db := config.DB
	var newArtefak model.Artefak

	// Ambil user_id dan role dari token
	userID, exists := c.Get("user_id")
	role, roleExists := c.Get("user_role")
	if !exists || !roleExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	// Hanya mahasiswa yang bisa membuat Artefak
	if role != "Mahasiswa" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Mahasiswa yang dapat mengajukan Artefak"})
		return
	}

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&newArtefak); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set user_id dari token
	newArtefak.UserID = userID.(uint)

	// Simpan ke DB
	if err := db.Create(&newArtefak).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Artefak berhasil ditambahkan", "data": newArtefak})
}

// Update request Artefak (hanya milik sendiri)
func UpdateArtefak(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	var artefak model.Artefak
	if err := db.Where("artefak_id = ?", id).First(&artefak).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artefak tidak ditemukan"})
		return
	}

	// Pastikan user yang login adalah pemiliknya
	userID, exists := c.Get("user_id")
	if !exists || artefak.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk mengubah artefak ini"})
		return
	}

	// Bind JSON
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update data
	if err := db.Model(&artefak).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artefak berhasil diperbarui", "data": artefak})
}

// Hapus request artefak
func DeleteArtefak(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	var artefak model.Artefak
	if err := db.Where("artefak_id = ?", id).First(&artefak).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artefak tidak ditemukan"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists || artefak.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk menghapus artefak ini"})
		return
	}

	if err := db.Delete(&artefak).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "artefak berhasil dihapus"})
}
