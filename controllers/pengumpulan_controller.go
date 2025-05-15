package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
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
        c.JSON(http.StatusNotFound, gin.H{"error": "Kelompok tidak ditemukan untuk user"})
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
        Preload("TahunMasuk").
        Preload("PengumpulanTugas", "kelompok_id = ?", km.KelompokID).
        Find(&tugasList).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    fmt.Printf("Found %d tugas\n", len(tugasList))

    // Format the response to match what the Flutter app expects
    var response []map[string]interface{}
    for _, tugas := range tugasList {
        // Check if there's a submission for this tugas by this kelompok
        var submission map[string]interface{} = nil
        if len(tugas.PengumpulanTugas) > 0 {
            pengumpulan := tugas.PengumpulanTugas[0]
            submission = map[string]interface{}{
                "id": pengumpulan.ID,
                "kelompok_id": pengumpulan.KelompokID,
                "tugas_id": pengumpulan.TugasID,
                "waktu_submit": pengumpulan.WaktuSubmit.Format(time.RFC3339),
                "file_path": pengumpulan.FilePath,
                "status": pengumpulan.Status,
            }
        }

        // Process file path for Laravel storage if needed
        filePath := tugas.File
        if strings.HasPrefix(filePath, "tugas_files/") {
            // This is a Laravel storage path, keep as is
        } else if !strings.HasPrefix(filePath, "http") && filePath != "" {
            // This is a local path, make it accessible via our file proxy
            filePath = "uploads/" + filePath
        }

        // Replace ternary operator with if-else statements
        var kategoriPA map[string]interface{}
        if tugas.KategoriPA.ID > 0 {
            kategoriPA = map[string]interface{}{
                "id": tugas.KategoriPA.ID,
                "kategori_pa": tugas.KategoriPA.KategoriPA,
            }
        } else {
            kategoriPA = nil
        }

        var tahunMasuk map[string]interface{}
        if tugas.TahunMasuk.ID > 0 {
            tahunMasuk = map[string]interface{}{
                "id": tugas.TahunMasuk.ID,
                "tahun_masuk": tugas.TahunMasuk.TahunMasuk,
                "status": tugas.TahunMasuk.Status,
            }
        } else {
            tahunMasuk = nil
        }

        var pengumpulanTugas []map[string]interface{}
        if submission != nil {
            pengumpulanTugas = []map[string]interface{}{submission}
        } else {
            pengumpulanTugas = []map[string]interface{}{}
        }

        item := map[string]interface{}{
            "id": tugas.ID,
            "judul_tugas": tugas.JudulTugas,
            "deskripsi_tugas": tugas.DeskripsiTugas,
            "tanggal_pengumpulan": tugas.TanggalPengumpulan.Format(time.RFC3339),
            "file": filePath,
            "user_id": tugas.UserID,
            "status": tugas.Status,
            "kategori_tugas": tugas.KategoriTugas,
            "kpa_id": tugas.KPAID,
            "prodi_id": tugas.ProdiID,
            "tm_id": tugas.TMID,
            "prodi": map[string]interface{}{
                "id": tugas.Prodi.ID,
                "nama_prodi": tugas.Prodi.NamaProdi,
            },
            "kategori_pa": kategoriPA,
            "tahun_masuk": tahunMasuk,
            "pengumpulan_tugas": pengumpulanTugas,
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
        Preload("TahunMasuk").
        Preload("PengumpulanTugas", "kelompok_id = ?", km.KelompokID).
        First(&tugas).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Tugas tidak ditemukan"})
        return
    }

    // Check if there's a submission for this tugas by this kelompok
    var submission map[string]interface{} = nil
    if len(tugas.PengumpulanTugas) > 0 {
        pengumpulan := tugas.PengumpulanTugas[0]
        submission = map[string]interface{}{
            "id": pengumpulan.ID,
            "kelompok_id": pengumpulan.KelompokID,
            "tugas_id": pengumpulan.TugasID,
            "waktu_submit": pengumpulan.WaktuSubmit.Format(time.RFC3339),
            "file_path": pengumpulan.FilePath,
            "status": pengumpulan.Status,
        }
    }

    // Process file path for Laravel storage if needed
    filePath := tugas.File
    if strings.HasPrefix(filePath, "tugas_files/") {
        // This is a Laravel storage path, keep as is
    } else if !strings.HasPrefix(filePath, "http") && filePath != "" {
        // This is a local path, make it accessible via our file proxy
        filePath = "uploads/" + filePath
    }

    // Replace ternary operators with if-else statements
    var kategoriPA map[string]interface{}
    if tugas.KategoriPA.ID > 0 {
        kategoriPA = map[string]interface{}{
            "id": tugas.KategoriPA.ID,
            "kategori_pa": tugas.KategoriPA.KategoriPA,
        }
    } else {
        kategoriPA = nil
    }

    var tahunMasuk map[string]interface{}
    if tugas.TahunMasuk.ID > 0 {
        tahunMasuk = map[string]interface{}{
            "id": tugas.TahunMasuk.ID,
            "tahun_masuk": tugas.TahunMasuk.TahunMasuk,
            "status": tugas.TahunMasuk.Status,
        }
    } else {
        tahunMasuk = nil
    }

    var pengumpulanTugas []map[string]interface{}
    if submission != nil {
        pengumpulanTugas = []map[string]interface{}{submission}
    } else {
        pengumpulanTugas = []map[string]interface{}{}
    }

    // Format the response
    response := map[string]interface{}{
        "id": tugas.ID,
        "judul_tugas": tugas.JudulTugas,
        "deskripsi_tugas": tugas.DeskripsiTugas,
        "tanggal_pengumpulan": tugas.TanggalPengumpulan.Format(time.RFC3339),
        "file": filePath,
        "user_id": tugas.UserID,
        "status": tugas.Status,
        "kategori_tugas": tugas.KategoriTugas,
        "kpa_id": tugas.KPAID,
        "prodi_id": tugas.ProdiID,
        "tm_id": tugas.TMID,
        "prodi": map[string]interface{}{
            "id": tugas.Prodi.ID,
            "nama_prodi": tugas.Prodi.NamaProdi,
        },
        "kategori_pa": kategoriPA,
        "tahun_masuk": tahunMasuk,
        "pengumpulan_tugas": pengumpulanTugas,
    }

    // Return the result
    c.JSON(http.StatusOK, gin.H{
        "message": "Detail tugas ditemukan",
        "data": response,
    })
}

// UpdateUploadFileTugas handles updating an already uploaded assignment file
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

    // Check if submission already exists
    var pengumpulan model.PengumpulanTugas
    result := db.Where("kelompok_id = ? AND tugas_id = ?", km.KelompokID, tugasID).First(&pengumpulan)
    
    // Generate unique filename
    timestamp := time.Now().Unix()
    ext := filepath.Ext(file.Filename)
    filename := fmt.Sprintf("tugas_%d_%d_%d%s", km.KelompokID, tugasID, timestamp, ext)
    filePath := filepath.Join("uploads", "tugas", filename)

    // Save file to disk
    if err := c.SaveUploadedFile(file, filePath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
        return
    }

    if result.Error != nil {
        // Create new submission if it doesn't exist
        newPengumpulan := model.PengumpulanTugas{
            KelompokID:  km.KelompokID,
            TugasID:     uint(tugasID),
            WaktuSubmit: time.Now(),
            FilePath:    filePath,
            Status:      "Submitted",
        }

        if err := db.Create(&newPengumpulan).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengumpulan tugas"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "message": "File tugas berhasil dikumpulkan",
            "data": newPengumpulan,
        })
    } else {
        // Update existing submission
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
}