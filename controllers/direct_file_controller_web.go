package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// Ganti path ini sesuai lokasi file kamu
const LaravelStoragePath = "https://vokasitera.d4trpl-itdel.id/storage/tugas_files/"

// Fungsi untuk akses file langsung (view di browser)
func DirectFileAccess(c *gin.Context) {
	filePath := c.Param("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File path is required"})
		return
	}

	// Hilangkan slash depan jika ada
	filePath = strings.TrimPrefix(filePath, "/")

	// Debug: cek path yang diminta
	fmt.Printf("Requested file path: %s\n", filePath)

	// Gabungkan path storage dengan filePath
	fullPath := filepath.Join(LaravelStoragePath, filePath)
	fmt.Printf("Accessing file at: %s\n", fullPath)

	// Cek file ada atau tidak
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "File not found",
			"path":  fullPath,
		})
		return
	}

	// Buka file
	file, err := os.Open(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to open file",
			"details": err.Error(),
			"path":    fullPath,
		})
		return
	}
	defer file.Close()

	// Info file
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get file info",
			"details": err.Error(),
		})
		return
	}

	// Tentukan content type
	extension := strings.ToLower(filepath.Ext(fullPath))
	contentType := "application/octet-stream"
	switch extension {
	case ".pdf":
		contentType = "application/pdf"
	case ".doc", ".docx":
		contentType = "application/msword"
	case ".zip":
		contentType = "application/zip"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	filename := filepath.Base(filePath)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))

	// Kirim file ke client
	c.Status(http.StatusOK)
	io.Copy(c.Writer, file)
}

// Fungsi untuk download file (force download)
func DirectFileDownload(c *gin.Context) {
	filePath := c.Param("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File path is required"})
		return
	}

	filePath = strings.TrimPrefix(filePath, "/")
	fmt.Printf("Requested file path for download: %s\n", filePath)
	fullPath := filepath.Join(LaravelStoragePath, filePath)
	fmt.Printf("Accessing file at: %s\n", fullPath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "File not found",
			"path":  fullPath,
		})
		return
	}

	file, err := os.Open(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to open file",
			"details": err.Error(),
			"path":    fullPath,
		})
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get file info",
			"details": err.Error(),
		})
		return
	}

	extension := strings.ToLower(filepath.Ext(fullPath))
	contentType := "application/octet-stream"
	switch extension {
	case ".pdf":
		contentType = "application/pdf"
	case ".doc", ".docx":
		contentType = "application/msword"
	case ".zip":
		contentType = "application/zip"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	filename := filepath.Base(filePath)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	c.Status(http.StatusOK)
	io.Copy(c.Writer, file)
}
