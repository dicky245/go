package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
	"fmt"
)

// GET - Ambil semua tugas yang sudah di-upload dari web
// Fungsi untuk mengambil semua tugas yang ditujukan ke kelompok ini (submitan dari dosen)

// GetSubmitanTugas retrieves assignments based on user's group
func GetSubmitanTugas(c *gin.Context) {
    db := config.DB

    // Get user_id from context
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan"})
        return
    }
    
    fmt.Printf("User ID: %v\n", userID)

    // Find kelompok mahasiswa based on user_id WITH PRELOADED KELOMPOK
    var km model.KelompokMahasiswa
    if err := db.Debug().
        Where("user_id = ?", userID).
        Preload("Kelompok").
        First(&km).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Kelompok tidak ditemukan untuk user"})
        return
    }
    
    fmt.Printf("Kelompok Mahasiswa: %+v\n", km)
    fmt.Printf("Kelompok: %+v\n", km.Kelompok)

    // Use the preloaded Kelompok
    kelompok := km.Kelompok
    
    // Log the query parameters for debugging
    fmt.Printf("Querying tugas with prodi_id = %d and TA_id = %d\n", kelompok.ProdiID, kelompok.TAID)

    // Query for tugas using the values from the user's kelompok
    var tugasList []model.Tugas
    if err := db.Debug().
        Where("prodi_id = ? AND TA_id = ?", kelompok.ProdiID, kelompok.TAID).
        Preload("Prodi").
        Preload("KategoriPA").
        Preload("TahunAjaran").
        Find(&tugasList).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    fmt.Printf("Found %d tugas\n", len(tugasList))

    // Format the response to match what the Flutter app expects
    var response []map[string]interface{}
    for _, tugas := range tugasList {
        item := map[string]interface{}{
            "id": tugas.ID,
            "judul": tugas.JudulTugas,
            "instruksi": tugas.DeskripsiTugas,
            "batas": tugas.TanggalPengumpulan.Format(time.RFC3339),
            "file": tugas.File,
            "userId": tugas.UserID,
            "status": tugas.Status,
            "prodi": tugas.Prodi.NamaProdi,
            "kategori_pa": tugas.KategoriPA.KategoriPA,
            "tahun_ajaran": tugas.TahunAjaran.TahunAjaran,
        }
        response = append(response, item)
    }

    // Return the results
    c.JSON(http.StatusOK, gin.H{
        "message": "Submitan tugas ditemukan",
        "data": response,
    })
}

// GetSubmitanTugasById retrieves a specific assignment by ID
func GetSubmitanTugasById(c *gin.Context) {
    db := config.DB
    
    // Get task ID from URL parameter
    tugasID := c.Param("id")
    
    // Get user_id from context
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan"})
        return
    }
    
    // Find kelompok mahasiswa based on user_id
    var km model.KelompokMahasiswa
    if err := db.Where("user_id = ?", userID).First(&km).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Kelompok tidak ditemukan untuk user"})
        return
    }
    
    // Get the specific tugas
    var tugas model.Tugas
    if err := db.
        Where("id = ?", tugasID).
        Preload("Prodi").
        Preload("KategoriPA").
        Preload("TahunAjaran").
        First(&tugas).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Tugas tidak ditemukan"})
        return
    }
    
    // Format the response
    response := map[string]interface{}{
        "id": tugas.ID,
        "judul": tugas.JudulTugas,
        "instruksi": tugas.DeskripsiTugas,
        "batas": tugas.TanggalPengumpulan.Format(time.RFC3339),
        "file": tugas.File,
        "userId": tugas.UserID,
        "status": tugas.Status,
        "prodi": tugas.Prodi.NamaProdi,
        "kategori_pa": tugas.KategoriPA.KategoriPA,
        "tahun_ajaran": tugas.TahunAjaran.TahunAjaran,
    }
    
    // Return the result
    c.JSON(http.StatusOK, gin.H{
        "message": "Detail tugas ditemukan",
        "data": response,
    })
}

// UploadFileTugas handles file uploads for assignments
// UpdateFileTugas handles updating an already uploaded assignment file
// UpdateUploadFileTugas handles updating an uploaded file for assignments
func UpdateUploadFileTugas(c *gin.Context) {
    db := config.DB

    // Get tugas ID from URL parameter
    tugasIDStr := c.Param("id")
    tugasID, err := strconv.ParseUint(tugasIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tugas ID"})
        return
    }

    // Get user_id from context
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan"})
        return
    }

    // Find kelompok mahasiswa based on user_id
    var km model.KelompokMahasiswa
    if err := db.Where("user_id = ?", userID).First(&km).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Kelompok tidak ditemukan untuk user"})
        return
    }

    // Get file from form
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan"})
        return
    }

    // Validate file size
    if file.Size > 100*1024*1024 { // 100MB limit
        c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file maksimal 100MB"})
        return
    }

    // Find existing submission
    var pengumpulan model.PengumpulanTugas
    if err := db.Where("kelompok_id = ? AND tugas_id = ?", km.KelompokID, tugasID).First(&pengumpulan).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Pengumpulan tugas tidak ditemukan, tidak bisa update"})
        return
    }

    // Generate new unique filename
    timestamp := time.Now().Unix()
    ext := filepath.Ext(file.Filename)
    filename := fmt.Sprintf("tugas_%d_%d_%d%s", km.KelompokID, tugasID, timestamp, ext)
    filePath := filepath.Join("uploads", "tugas", filename)

    // Save new file to disk
    if err := c.SaveUploadedFile(file, filePath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file baru"})
        return
    }

    // Update submission record
    pengumpulan.FilePath = filePath
    pengumpulan.WaktuSubmit = time.Now()
    pengumpulan.Status = "Resubmitted"

    if err := db.Save(&pengumpulan).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui pengumpulan tugas"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "File tugas berhasil diperbarui",
        "data": pengumpulan,
    })
}


