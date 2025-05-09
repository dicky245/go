package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
)

// GetPengumuman retrieves all active announcements
func GetPengumuman(c *gin.Context) {
	db := config.DB
	
	var pengumumans []model.Pengumuman
	
	// Get query parameters for filtering
	prodiID := c.DefaultQuery("prodi_id", "")
	kpaID := c.DefaultQuery("kpa_id", "")
	
	// Base query
	query := db.Where("status = ?", "aktif")
	
	// Apply filters if provided
	if prodiID != "" {
		query = query.Where("prodi_id = ?", prodiID)
	}
	
	if kpaID != "" {
		query = query.Where("KPA_id = ?", kpaID)
	}
	
	// Execute query with order by newest first
	if err := query.Order("tanggal_penulisan DESC").Find(&pengumumans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch announcements"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": pengumumans,
	})
}

// GetPengumumanByID retrieves a specific announcement by ID
func GetPengumumanByID(c *gin.Context) {
	db := config.DB
	
	id := c.Param("id")
	var pengumuman model.Pengumuman
	
	if err := db.First(&pengumuman, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Announcement not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": pengumuman,
	})
}

// CreatePengumuman creates a new announcement
func CreatePengumuman(c *gin.Context) {
	db := config.DB
	
	// Get user ID from token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Parse form data
	judul := c.PostForm("judul")
	deskripsi := c.PostForm("deskripsi")
	kpaID := c.PostForm("kpa_id")
	prodiID := c.PostForm("prodi_id")
	tmID := c.PostForm("tm_id")
	
	// Validate required fields
	if judul == "" || deskripsi == "" || kpaID == "" || prodiID == "" || tmID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}
	
	// Convert string IDs to uint
	kpaIDInt, _ := strconv.ParseUint(kpaID, 10, 64)
	prodiIDInt, _ := strconv.ParseUint(prodiID, 10, 64)
	tmIDInt, _ := strconv.ParseUint(tmID, 10, 64)
	
	// Create new announcement
	pengumuman := model.Pengumuman{
		Judul:            judul,
		Deskripsi:        deskripsi,
		TanggalPenulisan: time.Now(),
		Status:           "aktif",
		UserID:           userID.(uint),
		KPAID:            uint(kpaIDInt),
		ProdiID:          uint(prodiIDInt),
		TMID:             uint(tmIDInt),
	}
	
	// Handle file upload if present
	file, err := c.FormFile("file")
	if err == nil {
		filename := "uploads/pengumuman/" + strconv.FormatUint(uint64(time.Now().Unix()), 10) + "_" + file.Filename
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		pengumuman.File = filename
	}
	
	// Save to database
	if err := db.Create(&pengumuman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create announcement"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": pengumuman,
	})
}

// DeletePengumuman deletes an announcement
func DeletePengumuman(c *gin.Context) {
	db := config.DB
	
	id := c.Param("id")
	var pengumuman model	.Pengumuman
	
	// Check if announcement exists
	if err := db.First(&pengumuman, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Announcement not found"})
		return
	}
	
	// Delete the announcement
	if err := db.Delete(&pengumuman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete announcement"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Announcement deleted successfully",
	})
}