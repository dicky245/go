package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
)

// GetAllTugas retrieves all tugas
func GetAllTugas(c *gin.Context) {
	var tugas []model.Tugas
	
	// Get query parameters for filtering
	kpaID := c.Query("kpa_id")
	prodiID := c.Query("prodi_id")
	tmID := c.Query("tm_id")
	status := c.Query("status")
	kategori := c.Query("kategori_tugas")

	// Start building the query
	query := config.DB.Preload("Prodi").Preload("KategoriPA").Preload("TahunMasuk")

	// Apply filters if provided
	if kpaID != "" {
		query = query.Where("KPA_id = ?", kpaID)
	}
	if prodiID != "" {
		query = query.Where("prodi_id = ?", prodiID)
	}
	if tmID != "" {
		query = query.Where("TM_id = ?", tmID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if kategori != "" {
		query = query.Where("kategori_tugas = ?", kategori)
	}

	// Execute the query
	if err := query.Find(&tugas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data tugas",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data tugas berhasil diambil",
		"data":    tugas,
	})
}

// GetTugasById retrieves a specific tugas by ID
func GetTugasById(c *gin.Context) {
	id := c.Param("id")
	var tugas model.Tugas

	if err := config.DB.Preload("Prodi").Preload("KategoriPA").Preload("TahunMasuk").Preload("PengumpulanTugas").First(&tugas, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Tugas tidak ditemukan",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data tugas berhasil diambil",
		"data":    tugas,
	})
}

// CreateTugas creates a new tugas
func CreateTugas(c *gin.Context) {
	var input struct {
		UserID            uint      `json:"user_id" binding:"required"`
		JudulTugas        string    `json:"judul_tugas" binding:"required"`
		DeskripsiTugas    string    `json:"deskripsi_tugas" binding:"required"`
		KPAID             uint      `json:"kpa_id" binding:"required"`
		ProdiID           uint      `json:"prodi_id" binding:"required"`
		TMID              uint      `json:"tm_id" binding:"required"`
		TanggalPengumpulan string    `json:"tanggal_pengumpulan" binding:"required"`
		File              string    `json:"file"`
		Status            string    `json:"status"`
		KategoriTugas     string    `json:"kategori_tugas"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Parse the date string to time.Time
	tanggalPengumpulan, err := time.Parse(time.RFC3339, input.TanggalPengumpulan)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format tanggal tidak valid. Gunakan format ISO 8601 (YYYY-MM-DDTHH:MM:SSZ)",
			"error":   err.Error(),
		})
		return
	}

	// Create the tugas object
	tugas := model.Tugas{
		UserID:            input.UserID,
		JudulTugas:        input.JudulTugas,
		DeskripsiTugas:    input.DeskripsiTugas,
		KPAID:             input.KPAID,
		ProdiID:           input.ProdiID,
		TMID:              input.TMID,
		TanggalPengumpulan: tanggalPengumpulan,
		File:              input.File,
		Status:            input.Status,
		KategoriTugas:     input.KategoriTugas,
	}

	// Set default values if not provided
	if tugas.Status == "" {
		tugas.Status = "berlangsung"
	}
	if tugas.KategoriTugas == "" {
		tugas.KategoriTugas = "Tugas"
	}

	// Save to database
	if err := config.DB.Create(&tugas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat tugas",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Tugas berhasil dibuat",
		"data":    tugas,
	})
}

// UpdateTugas updates an existing tugas
func UpdateTugas(c *gin.Context) {
	id := c.Param("id")
	var tugas model.Tugas

	// Check if tugas exists
	if err := config.DB.First(&tugas, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Tugas tidak ditemukan",
			"error":   err.Error(),
		})
		return
	}

	// Bind input data
	var input struct {
		JudulTugas        string    `json:"judul_tugas"`
		DeskripsiTugas    string    `json:"deskripsi_tugas"`
		KPAID             uint      `json:"kpa_id"`
		ProdiID           uint      `json:"prodi_id"`
		TMID              uint      `json:"tm_id"`
		TanggalPengumpulan string    `json:"tanggal_pengumpulan"`
		File              string    `json:"file"`
		Status            string    `json:"status"`
		KategoriTugas     string    `json:"kategori_tugas"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Update fields if provided
	if input.JudulTugas != "" {
		tugas.JudulTugas = input.JudulTugas
	}
	if input.DeskripsiTugas != "" {
		tugas.DeskripsiTugas = input.DeskripsiTugas
	}
	if input.KPAID != 0 {
		tugas.KPAID = input.KPAID
	}
	if input.ProdiID != 0 {
		tugas.ProdiID = input.ProdiID
	}
	if input.TMID != 0 {
		tugas.TMID = input.TMID
	}
	if input.TanggalPengumpulan != "" {
		tanggalPengumpulan, err := time.Parse(time.RFC3339, input.TanggalPengumpulan)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Format tanggal tidak valid. Gunakan format ISO 8601 (YYYY-MM-DDTHH:MM:SSZ)",
				"error":   err.Error(),
			})
			return
		}
		tugas.TanggalPengumpulan = tanggalPengumpulan
	}
	if input.File != "" {
		tugas.File = input.File
	}
	if input.Status != "" {
		tugas.Status = input.Status
	}
	if input.KategoriTugas != "" {
		tugas.KategoriTugas = input.KategoriTugas
	}

	// Save updates to database
	if err := config.DB.Save(&tugas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui tugas",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Tugas berhasil diperbarui",
		"data":    tugas,
	})
}

// DeleteTugas deletes a tugas
func DeleteTugas(c *gin.Context) {
	id := c.Param("id")
	var tugas model.Tugas

	// Check if tugas exists
	if err := config.DB.First(&tugas, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Tugas tidak ditemukan",
			"error":   err.Error(),
		})
		return
	}

	// Delete from database
	if err := config.DB.Delete(&tugas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus tugas",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Tugas berhasil dihapus",
	})
}

// GetTugasByKategori retrieves tugas filtered by kategori_tugas
func GetTugasByKategori(c *gin.Context) {
	kategori := c.Param("kategori")
	var tugas []model.Tugas

	if err := config.DB.Where("kategori_tugas = ?", kategori).
		Preload("Prodi").Preload("KategoriPA").Preload("TahunMasuk").
		Find(&tugas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data tugas",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data tugas berhasil diambil",
		"data":    tugas,
	})
}

// GetTugasByStatus retrieves tugas filtered by status
func GetTugasByStatus(c *gin.Context) {
	status := c.Param("status")
	var tugas []model.Tugas

	if err := config.DB.Where("status = ?", status).
		Preload("Prodi").Preload("KategoriPA").Preload("TahunMasuk").
		Find(&tugas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data tugas",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data tugas berhasil diambil",
		"data":    tugas,
	})
}

// GetTugasByProdi retrieves tugas filtered by prodi_id
func GetTugasByProdi(c *gin.Context) {
	prodiID, err := strconv.Atoi(c.Param("prodi_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID Prodi tidak valid",
			"error":   err.Error(),
		})
		return
	}

	var tugas []model.Tugas

	if err := config.DB.Where("prodi_id = ?", prodiID).
		Preload("Prodi").Preload("KategoriPA").Preload("TahunMasuk").
		Find(&tugas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data tugas",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data tugas berhasil diambil",
		"data":    tugas,
	})
}