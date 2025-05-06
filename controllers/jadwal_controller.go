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

		// Get kelompok details to find prodi, KPA, and TM
		var kelompok []model.Kelompok
		if err := db.Where("id IN ?", kelompokIDs).Find(&kelompok).Error; err != nil {
			fmt.Printf("Error fetching kelompok details: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch kelompok details"})
			return
		}

		// Get all jadwal for the user's kelompok - MODIFIED QUERY
		var jadwal []struct {
			model.Jadwal
			KelompokNama string `json:"kelompok_nama"`
		}

		// Modified query without dosen table references
		query := `
			SELECT j.*, k.nomor_kelompok as kelompok_nama
			FROM jadwal j
			JOIN kelompok k ON j.kelompok_id = k.id
			WHERE j.kelompok_id IN (?)
			ORDER BY j.waktu DESC
		`

		if err := db.Raw(query, kelompokIDs).Scan(&jadwal).Error; err != nil {
			fmt.Printf("Error fetching jadwal: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jadwal"})
			return
		}

		fmt.Printf("Found %d jadwal for user %v\n", len(jadwal), userID)
		for i, j := range jadwal {
			fmt.Printf("Jadwal %d: ID=%v, KelompokID=%v, Ruangan=%v, Waktu=%v\n", 
				i, j.ID, j.KelompokID, j.Ruangan, j.Waktu)
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   jadwal,
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

		// Get the specific jadwal with additional information - MODIFIED QUERY
		var jadwal struct {
			model.Jadwal
			KelompokNama string `json:"kelompok_nama"`
		}

		// Modified query without dosen table references
		query := `
			SELECT j.*, k.nomor_kelompok as kelompok_nama
			FROM jadwal j
			JOIN kelompok k ON j.kelompok_id = k.id
			WHERE j.id = ? AND j.kelompok_id IN (?)
			LIMIT 1
		`

		if err := db.Raw(query, jadwalID, kelompokIDs).Scan(&jadwal).Error; err != nil {
			fmt.Printf("Error fetching jadwal details: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jadwal details"})
			return
		}

		if jadwal.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Jadwal not found or not authorized"})
			return
		}

		fmt.Printf("Found jadwal: ID=%v, KelompokID=%v, Ruangan=%v, Waktu=%v\n", 
			jadwal.ID, jadwal.KelompokID, jadwal.Ruangan, jadwal.Waktu)

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   jadwal,
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
			KelompokID uint      `json:"kelompok_id" binding:"required"`
			Ruangan    string    `json:"ruangan" binding:"required"`
			Waktu      time.Time `json:"waktu" binding:"required"`
			Penguji1   uint      `json:"penguji1" binding:"required"`
			Penguji2   uint      `json:"penguji2" binding:"required"`
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

		// Create the jadwal
		jadwal := model.Jadwal{
			KelompokID: request.KelompokID,
			Ruangan:    request.Ruangan,
			Waktu:      request.Waktu,
			UserID:     userID.(uint),
			Penguji1:   request.Penguji1,
			Penguji2:   request.Penguji2,
		}

		if err := db.Create(&jadwal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create jadwal"})
			return
		}

		// Send notification if needed
		// This is optional and can be implemented separately

		c.JSON(http.StatusCreated, gin.H{
			"status": "success",
			"data":   jadwal,
		})
	}