package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
)

// CREATE - Dosen buat pengumpulan
func CreatePengumpulan(c *gin.Context) {
	db := config.DB
	var newPengumpulan model.Submit

	userID, exists := c.Get("user_id")
	role, roleExists := c.Get("user_role")
	if !exists || !roleExists || role != "Dosen" {
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

// READ - Semua user bisa lihat pengumpulan
func GetPengumpulan(c *gin.Context) {
	db := config.DB
	userID, exists := c.Get("user_id")
	role, _ := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var pengumpulanList []model.Submit
	var err error

	if role == "Dosen" {
		err = db.Where("user_id = ?", userID).Find(&pengumpulanList).Error
	} else if role == "Mahasiswa" {
		err = db.Find(&pengumpulanList).Error
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pengumpulanList)
}

// READ BY ID - Semua user bisa lihat pengumpulan spesifik
func GetPengumpulanByID(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	var pengumpulan model.Submit
	if err := db.Where("submit_id = ?", id).First(&pengumpulan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengumpulan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, pengumpulan)
}

// UPDATE - Hanya dosen bisa update submitan miliknya
func UploadFilePengumpulan(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	var pengumpulan model.Submit
	if err := db.Where("submit_id = ?", id).First(&pengumpulan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengumpulan tidak ditemukan"})
		return
	}

	userID, exists := c.Get("user_id")
	role, _ := c.Get("user_role")
	if !exists || role != "Dosen" || pengumpulan.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk mengunggah file ini"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan"})
		return
	}

	if file.Size > 100*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file maksimal 100MB"})
		return
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	ext := filepath.Ext(file.Filename)
	safeFilename := "submit_" + id + "_" + timestamp + ext
	filePath := filepath.Join("uploads", safeFilename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pengumpulan.File = safeFilename
	if err := db.Save(&pengumpulan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File berhasil diunggah",
		"data":    pengumpulan,
	})
}

// DELETE - Hanya dosen bisa hapus submitan miliknya
func DeletePengumpulan(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	var pengumpulan model.Submit
	if err := db.Where("submit_id = ?", id).First(&pengumpulan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengumpulan tidak ditemukan"})
		return
	}

	userID, exists := c.Get("user_id")
	role, _ := c.Get("user_role")
	if !exists || role != "Dosen" || pengumpulan.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk menghapus submitan ini"})
		return
	}

	if err := db.Delete(&pengumpulan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Submitan berhasil dihapus"})
}
