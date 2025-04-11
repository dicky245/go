package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
)

func GetJadwal(c *gin.Context) {
	db := config.DB

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	var jadwalList []model.Jadwal
	if err := db.Where("user_id = ?", userID.(uint)).Find(&jadwalList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jadwalList)
}

func GetJadwalByID(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	role, _ := c.Get("user_role")
	if role == "Dosen" {
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
			return
		}

		var jadwal model.Jadwal
		if err := db.Where("jadwal_id = ? AND user_id = ?", id, userID).First(&jadwal).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Jadwal tidak ditemukan atau bukan milik Anda"})
			return
		}
		c.JSON(http.StatusOK, jadwal)
	} else if role == "Mahasiswa" {
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
			return
		}
		var jadwal model.Jadwal
		if err := db.Where("jadwal_id = ? AND kelompok = ?", id, userID).First(&jadwal).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Jadwal tidak ditemukan atau bukan milik Anda"})
			return
		}

		c.JSON(http.StatusOK, jadwal)
	}

}

func CreateJadwal(c *gin.Context) {
	db := config.DB
	var newJadwal model.Jadwal

	userID, exists := c.Get("user_id")
	role, roleExists := c.Get("user_role")
	if !exists || !roleExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	if role != "Dosen" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Dosen yang dapat mengajukan Jadwal"})
		return
	}

	if err := c.ShouldBindJSON(&newJadwal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newJadwal.UserID = userID.(uint)

	if err := db.Create(&newJadwal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Jadwal berhasil ditambahkan", "data": newJadwal})
}

func UpdateJadwal(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	var jadwal model.Jadwal
	if err := db.Where("jadwal_id = ?", id).First(&jadwal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jadwal tidak ditemukan"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists || jadwal.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk mengubah jadwal ini"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Model(&jadwal).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "jadwal berhasil diperbarui", "data": jadwal})
}

func DeleteJadwal(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	var jadwal model.Jadwal
	if err := db.Where("jadwal_id = ?", id).First(&jadwal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "jadwal tidak ditemukan"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists || jadwal.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk menghapus jadwal ini"})
		return
	}

	if err := db.Delete(&jadwal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "jadwal berhasil dihapus"})
}
