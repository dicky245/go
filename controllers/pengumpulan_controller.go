package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"fmt"
)

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
        // Return a more friendly error with a specific status code for "no group" case
        c.JSON(http.StatusOK, gin.H{
            "message": "Mahasiswa belum memiliki kelompok",
            "status": "no_group",
            "data": []map[string]interface{}{},
        })
        return
    }
    
    fmt.Printf("Kelompok Mahasiswa: %+v\n", km)
    fmt.Printf("Kelompok: %+v\n", km.Kelompok)

    // Use the preloaded Kelompok
    kelompok := km.Kelompok
    
    // Log the query parameters for debugging
    fmt.Printf("Querying tugas with prodi_id = %d and TM_id = %d\n", kelompok.ProdiID, kelompok.TMID)

    // Query for tugas using the values from the user's kelompok
    var tugasList []model.Tugas
    if err := db.Debug().
        Where("prodi_id = ? AND TM_id = ?", kelompok.ProdiID, kelompok.TMID).
        Preload("Prodi").
        Preload("KategoriPA").
        Find(&tugasList).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    fmt.Printf("Found %d tugas\n", len(tugasList))

    // Format the response to match what the Flutter app expects
    var response []map[string]interface{}
    for _, tugas := range tugasList {
        // Get tahun masuk data directly since the preload might fail
        var tahunMasuk model.TahunAjaran
        db.First(&tahunMasuk, tugas.TMID)
        
        // Check if this tugas has been submitted by the current kelompok
        var pengumpulan model.PengumpulanTugas
        var submissionStatus string = "Belum"
        var submissionFile string = ""
        
        if err := db.Where("kelompok_id = ? AND tugas_id = ?", km.KelompokID, tugas.ID).First(&pengumpulan).Error; err == nil {
            // Submission found
            submissionStatus = pengumpulan.Status
            submissionFile = pengumpulan.FilePath
        }
        
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
            "tahun_ajaran": tahunMasuk.TahunAjaran,
            "submission_status": submissionStatus,
            "submission_file": submissionFile,
        }
        response = append(response, item)
    }

    // Return the results
    c.JSON(http.StatusOK, gin.H{
        "message": "Submitan tugas ditemukan",
        "status": "success",
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
        // Return a more friendly error with a specific status code for "no group" case
        c.JSON(http.StatusOK, gin.H{
            "message": "Mahasiswa belum memiliki kelompok",
            "status": "no_group",
            "data": map[string]interface{}{},
        })
        return
    }
    
    // Get the specific tugas
    var tugas model.Tugas
    if err := db.
        Where("id = ?", tugasID).
        Preload("Prodi").
        Preload("KategoriPA").
        First(&tugas).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Tugas tidak ditemukan"})
        return
    }
    
    // Get tahun masuk data directly
    var tahunMasuk model.TahunAjaran
    db.First(&tahunMasuk, tugas.TMID)
    
    // Check if this tugas has been submitted by the current kelompok
    var pengumpulan model.PengumpulanTugas
    var submissionStatus string = "Belum"
    var submissionFile string = ""
    var submissionDate time.Time
    
    if err := db.Where("kelompok_id = ? AND tugas_id = ?", km.KelompokID, tugas.ID).First(&pengumpulan).Error; err == nil {
        // Submission found
        submissionStatus = pengumpulan.Status
        submissionFile = pengumpulan.FilePath
        submissionDate = pengumpulan.WaktuSubmit
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
        "tahun_ajaran": tahunMasuk.TahunAjaran,
        "submission_status": submissionStatus,
        "submission_file": submissionFile,
        "submission_date": submissionDate,
    }
    
    // Return the result
    c.JSON(http.StatusOK, gin.H{
        "message": "Detail tugas ditemukan",
        "status": "success",
        "data": response,
    })
}    

// UploadFileTugas handles file uploads for assignments
func UploadFileTugas(c *gin.Context) {
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

    // Check if submission already exists
    var existingSubmission model.PengumpulanTugas
    if err := db.Where("kelompok_id = ? AND tugas_id = ?", km.KelompokID, tugasID).First(&existingSubmission).Error; err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Tugas sudah pernah dikumpulkan. Gunakan fitur edit untuk memperbarui."})
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

    // Validate file extension
    ext := strings.ToLower(filepath.Ext(file.Filename))
    if ext != ".pdf" && ext != ".doc" && ext != ".docx" && ext != ".zip" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Format file tidak didukung. Hanya menerima file PDF, DOC, DOCX, atau ZIP"})
        return
    }

    // Generate unique filename
    timestamp := time.Now().Unix()
    filename := fmt.Sprintf("tugas_%d_%d_%d%s", km.KelompokID, tugasID, timestamp, ext)
    filePath := filepath.Join("uploads", "tugas", filename)

    // Save file to disk
    if err := c.SaveUploadedFile(file, filePath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
        return
    }

    // Check if deadline has passed
    var tugas model.Tugas
    if err := db.First(&tugas, tugasID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Tugas tidak ditemukan"})
        return
    }

    // Determine submission status
    submissionStatus := "Submitted"
    if time.Now().After(tugas.TanggalPengumpulan) {
        submissionStatus = "Late"
    }

    // Create submission record
    submission := model.PengumpulanTugas{
        KelompokID:  km.KelompokID,
        TugasID:     uint(tugasID),
        WaktuSubmit: time.Now(),
        FilePath:    filePath,
        Status:      submissionStatus,
    }

    if err := db.Create(&submission).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data pengumpulan"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "File tugas berhasil diunggah",
        "data": submission,
    })
}

// UpdateFileTugas handles updating an already uploaded assignment file
func UpdateFileTugas(c *gin.Context) {
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

    // Find existing submission
    var pengumpulan model.PengumpulanTugas
    if err := db.Where("kelompok_id = ? AND tugas_id = ?", km.KelompokID, tugasID).First(&pengumpulan).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Pengumpulan tugas tidak ditemukan, silakan submit terlebih dahulu"})
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

    // Validate file extension
    ext := strings.ToLower(filepath.Ext(file.Filename))
    if ext != ".pdf" && ext != ".doc" && ext != ".docx" && ext != ".zip" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Format file tidak didukung. Hanya menerima file PDF, DOC, DOCX, atau ZIP"})
        return
    }

    // Generate new unique filename
    timestamp := time.Now().Unix()
    filename := fmt.Sprintf("tugas_%d_%d_%d%s", km.KelompokID, tugasID, timestamp, ext)
    filePath := filepath.Join("uploads", "tugas", filename)

    // Save new file to disk
    if err := c.SaveUploadedFile(file, filePath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file baru"})
        return
    }

    // Check if deadline has passed
    var tugas model.Tugas
    if err := db.First(&tugas, tugasID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Tugas tidak ditemukan"})
        return
    }

    // Determine submission status
    submissionStatus := "Resubmitted"
    if time.Now().After(tugas.TanggalPengumpulan) && pengumpulan.Status != "Late" {
        submissionStatus = "Late"
    } else if pengumpulan.Status == "Late" {
        submissionStatus = "Late"
    }

    // Update submission record
    pengumpulan.FilePath = filePath
    pengumpulan.WaktuSubmit = time.Now()
    pengumpulan.Status = submissionStatus

    if err := db.Save(&pengumpulan).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui pengumpulan tugas"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "File tugas berhasil diperbarui",
        "data": pengumpulan,
    })
}