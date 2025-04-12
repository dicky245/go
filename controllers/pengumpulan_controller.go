package controllers

// controller submit
import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
)

func GetPengumpulan(c *gin.Context) {
	db := config.DB

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	var pengumpulanList []model.Submit
	if err := db.Where("user_id = ?", userID.(uint)).Find(&pengumpulanList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, pengumpulanList)
}

func GetPengumpulanByID(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	var pengumpulan model.Submit
	if err := db.Where("submit_id = ? AND user_id = ?", id, userID).First(&pengumpulan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengumpulan tidak ditemukan atau bukan milik Anda"})
		return
	}

	c.JSON(http.StatusOK, pengumpulan)
}

func CreatePengumpulan(c *gin.Context) {
	db := config.DB
	var newPengumpulan model.Submit

	userID, exists := c.Get("user_id")
	role, roleExists := c.Get("user_role")
	if !exists || !roleExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	if role != "Dosen" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Dosen yang dapat memberi pengumpulan"})
		return
	}

	if err := c.ShouldBindJSON(&newPengumpulan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPengumpulan.UserID = userID.(uint)

	if err := db.Create(&newPengumpulan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pengumpulan berhasil ditambahkan", "data": newPengumpulan})
}

func UpdatePengumpulan(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	var pengumpulan model.Submit
	if err := db.Where("submit_id = ?", id).First(&pengumpulan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengumpulan tidak ditemukan"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists || pengumpulan.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk mengubah artefak ini"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Model(&pengumpulan).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artefak berhasil diperbarui", "data": pengumpulan})
}

func DeletePengumpulan(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	var pengumpulan model.Submit
	if err := db.Where("submit_id = ?", id).First(&pengumpulan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengumpulan tidak ditemukan"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists || pengumpulan.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk menghapus pengumpulan ini"})
		return
	}

	if err := db.Delete(&pengumpulan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "pengumpulan berhasil dihapus"})
}
