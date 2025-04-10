package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
	"gorm.io/gorm"
)

func CreateKelompok(c *gin.Context) {
	var kelompok model.Kelompok

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&kelompok); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi duplikat kombinasi nomor, jenis_pa, prodi, dan angkatan
	var existing model.Kelompok
	err := config.DB.Where("nomor = ? AND jenis_pa = ? AND prodi = ? AND angkatan = ?",
		kelompok.Nomor, kelompok.JenisPA, kelompok.Prodi, kelompok.Angkatan).First(&existing).Error

	if err == nil {
		// Ditemukan ⇒ Duplikat
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kelompok dengan nomor tersebut sudah ada untuk kombinasi jenis PA, prodi, dan angkatan"})
		return
	} else if err != gorm.ErrRecordNotFound {
		// Error lain (misalnya DB error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Simpan ke database
	if err := config.DB.Create(&kelompok).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, kelompok)
}

func GetKelompok(c *gin.Context) {
	var kelompoks []model.Kelompok
	config.DB.Find(&kelompoks)
	c.JSON(http.StatusOK, kelompoks)
}

func GetKelompokByID(c* gin.Context){
	var kelompok model.Kelompok

	id, err := strconv.ParseUint(c.Param("id"),10,32)
	if err !=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    if err := config.DB.First(&kelompok, uint(id)).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Kelompok not found"})
        return
    }
    c.JSON(http.StatusOK, kelompok)
}

// Update Role
func UpdateKelompok(c *gin.Context) {
    var kelompok model.Kelompok

    // Konversi id dari string ke uint
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    if err := config.DB.First(&kelompok, uint(id)).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Kelompok not found"})
        return
    }

    if err := c.ShouldBindJSON(&kelompok); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    kelompok.ID = uint(id) // Pastikan ID tidak berubah
    config.DB.Save(&kelompok)
    c.JSON(http.StatusOK, kelompok)
}

// Delete Role
func DeleteKelompok(c *gin.Context) {
    var kelompok model.Kelompok

    // Konversi id dari string ke uint
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    if err := config.DB.First(&kelompok, uint(id)).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
        return
    }

    config.DB.Delete(&kelompok)
    c.JSON(http.StatusOK, gin.H{"message": "Role deleted"})
}