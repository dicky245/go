package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
)

// CreatePengumpulanTugas handles the task submission by mahasiswa
func CreatePengumpulanTugas(c *gin.Context) {
	log.Println("CreatePengumpulanTugas called")
	
	// Get the tugas_id from the URL parameter
	tugasIDStr := c.Param("id")
	log.Printf("Received tugas_id: %s", tugasIDStr)
	
	tugasID, err := strconv.ParseUint(tugasIDStr, 10, 32)
	if err != nil {
		log.Printf("Error parsing tugas_id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID tugas tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Check if the tugas exists
	var tugas model.Tugas
	if err := config.DB.First(&tugas, tugasID).Error; err != nil {
		log.Printf("Tugas not found with ID %d: %v", tugasID, err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Tugas tidak ditemukan",
			"error":   err.Error(),
		})
		return
	}
	log.Printf("Found tugas: %s", tugas.JudulTugas)

	// Get the file from the form
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error getting file from form: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "File tidak ditemukan dalam request",
			"error":   err.Error(),
		})
		return
	}
	log.Printf("Received file: %s, size: %d", file.Filename, file.Size)

	// Validate file type
	allowedExtensions := []string{".pdf", ".doc", ".docx", ".zip"}
	ext := filepath.Ext(file.Filename)
	validExt := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			validExt = true
			break
		}
	}

	if !validExt {
		log.Printf("Invalid file extension: %s", ext)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format file tidak didukung. Hanya menerima file PDF, DOC, DOCX, atau ZIP",
		})
		return
	}

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Printf("Error creating uploads directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat direktori untuk menyimpan file",
			"error":   err.Error(),
		})
		return
	}

	// Create a unique filename to prevent overwriting
	filename := fmt.Sprintf("%d_%s_%s", tugasID, time.Now().Format("20060102150405"), file.Filename)
	filePath := filepath.Join("uploads", filename)
	log.Printf("File will be saved to: %s", filePath)

	// Save the file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("Error saving file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan file",
			"error":   err.Error(),
		})
		return
	}
	log.Printf("File saved successfully")

	// Get kelompok_id from the authenticated user or set a default value
	// In a real application, this should come from the authenticated user's session
	var kelompokID uint = 1 // Default value
	log.Printf("Using kelompok_id: %d", kelompokID)

	// Check if a submission already exists for this tugas and kelompok
	var existingSubmission model.PengumpulanTugas
	result := config.DB.Where("tugas_id = ? AND kelompok_id = ?", tugasID, kelompokID).First(&existingSubmission)

	if result.Error == nil {
		log.Printf("Found existing submission, updating it")
		// Submission already exists, update it instead of creating a new one
		existingSubmission.FilePath = filePath
		existingSubmission.WaktuSubmit = time.Now()
		existingSubmission.Status = "Resubmitted"

		if err := config.DB.Save(&existingSubmission).Error; err != nil {
			log.Printf("Error updating existing submission: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal memperbarui data pengumpulan tugas",
				"error":   err.Error(),
			})
			return
		}

		log.Printf("Submission updated successfully")
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Tugas berhasil diperbarui",
			"data":    existingSubmission,
		})
		return
	}

	// Create a new submission record
	pengumpulanTugas := model.PengumpulanTugas{
		TugasID:     uint(tugasID),
		KelompokID:  kelompokID,
		WaktuSubmit: time.Now(),
		FilePath:    filePath,
		Status:      "Submitted", // Initial status
	}
	log.Printf("Creating new submission: %+v", pengumpulanTugas)

	// Save to database
	if err := config.DB.Create(&pengumpulanTugas).Error; err != nil {
		log.Printf("Error creating new submission: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan data pengumpulan tugas",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("New submission created successfully with ID: %d", pengumpulanTugas.ID)
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Tugas berhasil dikumpulkan",
		"data":    pengumpulanTugas,
	})
}

// UpdateTugas handles updating an existing submission
func UpdateTugas(c *gin.Context) {
	log.Println("UpdateTugas called")
	
	// Get the submission ID from the URL parameter
	id := c.Param("id")
	log.Printf("Received tugas_id: %s", id)
	
	// Get the file from the form
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error getting file from form: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "File tidak ditemukan dalam request",
			"error":   err.Error(),
		})
		return
	}
	log.Printf("Received file: %s, size: %d", file.Filename, file.Size)

	// Validate file type
	allowedExtensions := []string{".pdf", ".doc", ".docx", ".zip"}
	ext := filepath.Ext(file.Filename)
	validExt := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			validExt = true
			break
		}
	}

	if !validExt {
		log.Printf("Invalid file extension: %s", ext)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format file tidak didukung. Hanya menerima file PDF, DOC, DOCX, atau ZIP",
		})
		return
	}

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Printf("Error creating uploads directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat direktori untuk menyimpan file",
			"error":   err.Error(),
		})
		return
	}

	// Get kelompok_id from the authenticated user or set a default value
	var kelompokID uint = 1 // Default value
	log.Printf("Using kelompok_id: %d", kelompokID)

	// Find the existing submission
	var pengumpulan model.PengumpulanTugas
	if err := config.DB.Where("tugas_id = ? AND kelompok_id = ?", id, kelompokID).First(&pengumpulan).Error; err != nil {
		log.Printf("No existing submission found, creating new one")
		// If no submission exists, create a new one
		tugasID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			log.Printf("Error parsing tugas_id: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "ID tugas tidak valid",
				"error":   err.Error(),
			})
			return
		}

		// Create a unique filename to prevent overwriting
		filename := fmt.Sprintf("%d_%s_%s", tugasID, time.Now().Format("20060102150405"), file.Filename)
		filePath := filepath.Join("uploads", filename)
		log.Printf("File will be saved to: %s", filePath)

		// Save the file
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			log.Printf("Error saving file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal menyimpan file",
				"error":   err.Error(),
			})
			return
		}
		log.Printf("File saved successfully")

		// Create a new submission
		newPengumpulan := model.PengumpulanTugas{
			TugasID:     uint(tugasID),
			KelompokID:  kelompokID,
			WaktuSubmit: time.Now(),
			FilePath:    filePath,
			Status:      "Submitted",
		}
		log.Printf("Creating new submission: %+v", newPengumpulan)

		if err := config.DB.Create(&newPengumpulan).Error; err != nil {
			log.Printf("Error creating new submission: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal menyimpan data pengumpulan tugas",
				"error":   err.Error(),
			})
			return
		}

		log.Printf("New submission created successfully with ID: %d", newPengumpulan.ID)
		c.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "Tugas berhasil dikumpulkan",
			"data":    newPengumpulan,
		})
		return
	}

	// Create a unique filename to prevent overwriting
	filename := fmt.Sprintf("%d_%s_%s", pengumpulan.TugasID, time.Now().Format("20060102150405"), file.Filename)
	filePath := filepath.Join("uploads", filename)
	log.Printf("File will be saved to: %s", filePath)

	// Save the file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("Error saving file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan file",
			"error":   err.Error(),
		})
		return
	}
	log.Printf("File saved successfully")

	// Update the submission record
	pengumpulan.FilePath = filePath
	pengumpulan.WaktuSubmit = time.Now()
	pengumpulan.Status = "Resubmitted" // Update status to indicate it's been resubmitted
	log.Printf("Updating submission: %+v", pengumpulan)

	// Save to database
	if err := config.DB.Save(&pengumpulan).Error; err != nil {
		log.Printf("Error updating submission: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui data pengumpulan tugas",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Submission updated successfully")
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pengumpulan tugas berhasil diperbarui",
		"data":    pengumpulan,
	})
}

// GetAllTugas retrieves all tugas assigned by dosen for mahasiswa
func GetAllTugas(c *gin.Context) {
	log.Println("GetAllTugas called")
	
	var tugas []model.Tugas
	
	// Get query parameters for filtering
	kpaID := c.Query("kpa_id")
	prodiID := c.Query("prodi_id")
	tmID := c.Query("tm_id")
	status := c.Query("status")
	kategori := c.Query("kategori_tugas")
	
	log.Printf("Filters - kpaID: %s, prodiID: %s, tmID: %s, status: %s, kategori: %s", 
		kpaID, prodiID, tmID, status, kategori)

	// Start building the query
	query := config.DB.Model(&model.Tugas{}).Preload("Prodi").Preload("KategoriPA").Preload("TahunMasuk").Preload("PengumpulanTugas")

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
		log.Printf("Error fetching tugas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data tugas",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Successfully fetched %d tugas", len(tugas))
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data tugas berhasil diambil",
		"data":    tugas,
	})
}

// GetTugasById retrieves a specific tugas by ID assigned by dosen
func GetTugasById(c *gin.Context) {
	id := c.Param("id")
	log.Printf("GetTugasById called with id: %s", id)
	
	var tugas model.Tugas

	// Fetch the tugas made by dosen and related to mahasiswa
	if err := config.DB.Preload("Prodi").Preload("KategoriPA").Preload("TahunMasuk").Preload("PengumpulanTugas").First(&tugas, id).Error; err != nil {
		log.Printf("Error fetching tugas with id %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Tugas tidak ditemukan",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Successfully fetched tugas with id %s: %s", id, tugas.JudulTugas)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data tugas berhasil diambil",
		"data":    tugas,
	})
}

// GetTugasByKategori retrieves tugas filtered by kategori_tugas
func GetTugasByKategori(c *gin.Context) {
	kategori := c.Param("kategori")
	log.Printf("GetTugasByKategori called with kategori: %s", kategori)
	
	var tugas []model.Tugas

	// Make sure we're using the exact case as defined in the enum
	// The database has 'Tugas', 'Revisi', 'Artefak' as enum values
	if err := config.DB.Where("kategori_tugas = ?", kategori).
		Preload("Prodi").Preload("KategoriPA").Preload("TahunMasuk").Preload("PengumpulanTugas").
		Find(&tugas).Error; err != nil {
		log.Printf("Error fetching tugas by kategori %s: %v", kategori, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data tugas",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Successfully fetched %d tugas with kategori %s", len(tugas), kategori)
	// Log the query results for debugging
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data tugas berhasil diambil",
		"data":    tugas,
		"count":   len(tugas),
		"filter":  kategori,
	})
}

// GetTugasByStatus retrieves tugas filtered by status
func GetTugasByStatus(c *gin.Context) {
	status := c.Param("status")
	log.Printf("GetTugasByStatus called with status: %s", status)
	
	var tugas []model.Tugas

	if err := config.DB.Where("status = ?", status).
		Preload("Prodi").Preload("KategoriPA").Preload("TahunMasuk").Preload("PengumpulanTugas").
		Find(&tugas).Error; err != nil {
		log.Printf("Error fetching tugas by status %s: %v", status, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data tugas",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Successfully fetched %d tugas with status %s", len(tugas), status)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data tugas berhasil diambil",
		"data":    tugas,
	})
}

// GetTugasByProdi retrieves tugas filtered by prodi_id
func GetTugasByProdi(c *gin.Context) {
	prodiID := c.Param("prodi_id")
	log.Printf("GetTugasByProdi called with prodi_id: %s", prodiID)
	
	var tugas []model.Tugas

	if err := config.DB.Where("prodi_id = ?", prodiID).
		Preload("Prodi").Preload("KategoriPA").Preload("TahunMasuk").Preload("PengumpulanTugas").
		Find(&tugas).Error; err != nil {
		log.Printf("Error fetching tugas by prodi_id %s: %v", prodiID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data tugas",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Successfully fetched %d tugas with prodi_id %s", len(tugas), prodiID)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data tugas berhasil diambil",
		"data":    tugas,
	})
}
