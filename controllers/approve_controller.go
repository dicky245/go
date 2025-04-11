package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
)

func GetUpdateBimbingan(c *gin.Context) {
	db := config.DB

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	role, _ := c.Get("user_role")
	var bimbinganList []model.Bimbingan

	if role == "Dosen" {
		if err := db.Find(&bimbinganList).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := db.Where("user_id = ?", userID.(uint)).Find(&bimbinganList).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, bimbinganList)
}

func GetUpdateBimbinganByID(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	role, _ := c.Get("user_role")

	var bimbingan model.Bimbingan
	if err := db.Where("bimbingan_id = ?", id).First(&bimbingan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bimbingan tidak ditemukan"})
		return
	}

	if role != "Dosen" && bimbingan.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke bimbingan ini"})
		return
	}

	c.JSON(http.StatusOK, bimbingan)
}

func UpdateRequestBimbingan(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	role, _ := c.Get("user_role")
	if role != "Dosen" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya dosen yang bisa mengubah status"})
		return
	}

	var bimbingan model.Bimbingan
	if err := db.Where("bimbingan_id = ?", id).First(&bimbingan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bimbingan tidak ditemukan"})
		return
	}

	var request struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Model(&bimbingan).Update("status", request.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status diperbarui", "data": bimbingan})
}
