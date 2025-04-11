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
		return
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

	userID, exists := c.Get("user_id")
	role, roleExists := c.Get("user_role")
	if !exists || !roleExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	if role != "Mahasiswa" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Mahasiswa yang dapat mengumpulkan Artefak"})
		return
	}

	if err := c.ShouldBindJSON(&newArtefak); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newArtefak.UserID = userID.(uint)

	if err := db.Create(&newArtefak).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Artefak berhasil ditambahkan", "data": newArtefak})
}

func UpdateArtefak(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	var artefak model.Artefak
	if err := db.Where("artefak_id = ?", id).First(&artefak).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artefak tidak ditemukan"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists || artefak.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk mengubah artefak ini"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Model(&artefak).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artefak berhasil diperbarui", "data": artefak})
}

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

	c.JSON(http.StatusOK, gin.H{"message": "Artefak berhasil dihapus"})
}

func GetSubmitMahasiswa(c *gin.Context) {
	db := config.DB
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan dalam token"})
		return
	}

	var submits []model.Submit
	if err := db.Find(&submits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var artefaks []model.Artefak
	if err := db.Where("user_id = ?", userID.(uint)).Find(&artefaks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	artefakMap := make(map[uint]model.Artefak)
	for _, a := range artefaks {
		artefakMap[a.SubmitID] = a
	}
	var response []gin.H
	for _, s := range submits {
		item := gin.H{
			"submit_id":   s.SubmitID,
			"judul":       s.Judul,
			"instruksi":   s.Instruksi,
			"deadline":    s.Batas,
			"is_uploaded": false,
			"artefak":     nil,
		}

		if a, found := artefakMap[s.SubmitID]; found {
			item["is_uploaded"] = true
			item["artefak"] = a
		}

		response = append(response, item)
	}

	c.JSON(http.StatusOK, response)
}
