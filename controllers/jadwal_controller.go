package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
)

// GetJadwal retrieves jadwal for the authenticated user based on their kelompok
func GetJadwal(c *gin.Context) {
	// Get user ID from token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	fmt.Printf("GetJadwal called for user ID: %v\n", userID)

	// Find the user's kelompok
	var kelompokMahasiswa []model.KelompokMahasiswa
	db := config.DB

	if err := db.Where("user_id = ?", userID).Find(&kelompokMahasiswa).Error; err != nil {
		fmt.Printf("Error fetching kelompok for user %v: %v\n", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user's kelompok"})
		return
	}

	fmt.Printf("Found %d kelompok for user %v\n", len(kelompokMahasiswa), userID)

	if len(kelompokMahasiswa) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   []interface{}{},
			"message": "User is not assigned to any kelompok",
		})
		return
	}

	// Get kelompok IDs
	var kelompokIDs []uint
	for _, km := range kelompokMahasiswa {
		kelompokIDs = append(kelompokIDs, km.KelompokID)
		fmt.Printf("User %v is in kelompok %v\n", userID, km.KelompokID)
	}

	// Get all jadwal for the user's kelompok with additional information
	var jadwalResults []struct {
		model.Jadwal
		KelompokNama string `json:"kelompok_nama"`
		Ruangan      string `json:"ruangan"`
	}

	// Query to get jadwal with kelompok name and ruangan name
	query := `
		SELECT j.*, k.nomor_kelompok as kelompok_nama, r.ruangan as ruangan
		FROM jadwal j
		JOIN kelompok k ON j.kelompok_id = k.id
		JOIN ruangan r ON j.ruangan_id = r.id
		WHERE j.kelompok_id IN (?)
		ORDER BY j.waktu_mulai DESC
	`

	if err := db.Raw(query, kelompokIDs).Scan(&jadwalResults).Error; err != nil {
		fmt.Printf("Error fetching jadwal: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jadwal"})
		return
	}

	// For each jadwal, find the penguji (examiners)
	var jadwalResponse []gin.H
	for _, j := range jadwalResults {
		// Find penguji for this kelompok
		var pengujiList []model.Penguji
		if err := db.Where("kelompok_id = ?", j.KelompokID).Find(&pengujiList).Error; err != nil {
			fmt.Printf("Error fetching penguji for kelompok %v: %v\n", j.KelompokID, err)
			// Continue without penguji info
		}

		// Create response with penguji info
		jadwalItem := gin.H{
			"id":            j.ID,
			"kelompok_id":   j.KelompokID,
			"ruangan":       j.Ruangan,
			"waktu":         j.WaktuMulai, // Use waktu_mulai as the main time field for compatibility
			"waktu_mulai":   j.WaktuMulai,
			"waktu_selesai": j.WaktuSelesai,
			"user_id":       j.UserID,
			"ruangan_id":    j.RuanganID,
			"kelompok_nama": j.KelompokNama,
			"created_at":    j.CreatedAt,
			"updated_at":    j.UpdatedAt,
		}

		// Add penguji info if available
		if len(pengujiList) > 0 {
			// Assume first penguji is penguji1 and second is penguji2
			jadwalItem["penguji1"] = pengujiList[0].UserID
			if len(pengujiList) > 1 {
				jadwalItem["penguji2"] = pengujiList[1].UserID
			} else {
				jadwalItem["penguji2"] = 0 // Default value if no second penguji
			}
		} else {
			jadwalItem["penguji1"] = 0
			jadwalItem["penguji2"] = 0
		}

		jadwalResponse = append(jadwalResponse, jadwalItem)
	}

	fmt.Printf("Found %d jadwal for user %v\n", len(jadwalResponse), userID)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   jadwalResponse,
	})
}

// GetJadwalByID retrieves a specific jadwal by ID
func GetJadwalByID(c *gin.Context) {
	// Get user ID from token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	jadwalID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid jadwal ID"})
		return
	}

	fmt.Printf("GetJadwalByID called for user %v, jadwal ID %v\n", userID, jadwalID)

	// Find the user's kelompok
	var kelompokMahasiswa []model.KelompokMahasiswa
	db := config.DB

	if err := db.Where("user_id = ?", userID).Find(&kelompokMahasiswa).Error; err != nil {
		fmt.Printf("Error fetching kelompok for user %v: %v\n", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user's kelompok"})
		return
	}

	// Get kelompok IDs
	var kelompokIDs []uint
	for _, km := range kelompokMahasiswa {
		kelompokIDs = append(kelompokIDs, km.KelompokID)
	}

	// Get the specific jadwal with additional information
	var jadwalResult struct {
		model.Jadwal
		KelompokNama string `json:"kelompok_nama"`
		Ruangan      string `json:"ruangan"`
	}

	// Query to get jadwal with kelompok name and ruangan name
	query := `
		SELECT j.*, k.nomor_kelompok as kelompok_nama, r.ruangan as ruangan
		FROM jadwal j
		JOIN kelompok k ON j.kelompok_id = k.id
		JOIN ruangan r ON j.ruangan_id = r.id
		WHERE j.id = ? AND j.kelompok_id IN (?)
		LIMIT 1
	`

	if err := db.Raw(query, jadwalID, kelompokIDs).Scan(&jadwalResult).Error; err != nil {
		fmt.Printf("Error fetching jadwal details: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jadwal details"})
		return
	}

	if jadwalResult.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jadwal not found or not authorized"})
		return
	}

	// Find penguji for this kelompok
	var pengujiList []model.Penguji
	if err := db.Where("kelompok_id = ?", jadwalResult.KelompokID).Find(&pengujiList).Error; err != nil {
		fmt.Printf("Error fetching penguji for kelompok %v: %v\n", jadwalResult.KelompokID, err)
		// Continue without penguji info
	}

	// Create response with penguji info
	jadwalResponse := gin.H{
		"id":            jadwalResult.ID,
		"kelompok_id":   jadwalResult.KelompokID,
		"ruangan":       jadwalResult.Ruangan,
		"waktu":         jadwalResult.WaktuMulai, // Use waktu_mulai as the main time field for compatibility
		"waktu_mulai":   jadwalResult.WaktuMulai,
		"waktu_selesai": jadwalResult.WaktuSelesai,
		"user_id":       jadwalResult.UserID,
		"ruangan_id":    jadwalResult.RuanganID,
		"kelompok_nama": jadwalResult.KelompokNama,
		"created_at":    jadwalResult.CreatedAt,
		"updated_at":    jadwalResult.UpdatedAt,
	}

	// Add penguji info if available
	if len(pengujiList) > 0 {
		// Assume first penguji is penguji1 and second is penguji2
		jadwalResponse["penguji1"] = pengujiList[0].UserID
		if len(pengujiList) > 1 {
			jadwalResponse["penguji2"] = pengujiList[1].UserID
		} else {
			jadwalResponse["penguji2"] = 0 // Default value if no second penguji
		}
	} else {
		jadwalResponse["penguji1"] = 0
		jadwalResponse["penguji2"] = 0
	}

	fmt.Printf("Found jadwal: ID=%v, KelompokID=%v, Ruangan=%v, Waktu=%v\n", 
		jadwalResult.ID, jadwalResult.KelompokID, jadwalResult.Ruangan, jadwalResult.WaktuMulai)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   jadwalResponse,
	})
}

// GetRuangan retrieves all available ruangan (rooms)
func GetRuangan(c *gin.Context) {
	var ruangan []model.Ruangan
	db := config.DB

	if err := db.Order("ruangan").Find(&ruangan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ruangan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   ruangan,
	})
}

// CreateJadwal creates a new jadwal (for admin/dosen)
func CreateJadwal(c *gin.Context) {
	// Get user ID from token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var request struct {
		KelompokID   uint      `json:"kelompok_id" binding:"required"`
		RuanganID    uint      `json:"ruangan_id" binding:"required"`
		WaktuMulai   time.Time `json:"waktu_mulai" binding:"required"`
		WaktuSelesai time.Time `json:"waktu_selesai" binding:"required"`
		KPAID        uint      `json:"kpa_id" binding:"required"`
		ProdiID      uint      `json:"prodi_id" binding:"required"`
		TMID         uint      `json:"tm_id" binding:"required"`
		Penguji      []uint    `json:"penguji" binding:"required,min=1"` // List of penguji user IDs
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the kelompok exists
	var kelompok model.Kelompok
	db := config.DB
	if err := db.First(&kelompok, request.KelompokID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kelompok not found"})
		return
	}

	// Check if the ruangan exists
	var ruangan model.Ruangan
	if err := db.First(&ruangan, request.RuanganID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ruangan not found"})
		return
	}

	// Create the jadwal
	jadwal := model.Jadwal{
		KelompokID:   request.KelompokID,
		RuanganID:    request.RuanganID,
		WaktuMulai:   request.WaktuMulai,
		WaktuSelesai: request.WaktuSelesai,
		UserID:       userID.(uint),
		KPAID:        request.KPAID,
		ProdiID:      request.ProdiID,
		TMID:         request.TMID,
	}

	// Start a transaction
	tx := db.Begin()

	// Create jadwal
	if err := tx.Create(&jadwal).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create jadwal"})
		return
	}

	// Create penguji records
	for i, pengujiUserID := range request.Penguji {
		if i >= 2 {
			break // Only support up to 2 penguji
		}

		penguji := model.Penguji{
			UserID:     pengujiUserID,
			KelompokID: request.KelompokID,
		}

		if err := tx.Create(&penguji).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create penguji"})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Send notification if needed
	// This is optional and can be implemented separately

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   jadwal,
	})
}