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

// Configuration for direct file access
const (
	// Update this to the actual path to your Laravel storage directory
	LaravelStoragePath = "D:/Semester 4/PA II/vokasiteraweb/storage/app/public/"
)

// DirectFileAccess serves files directly from the filesystem
func DirectFileAccess(c *gin.Context) {
	// Get the file path from the URL
	filePath := c.Param("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File path is required"})
		return
	}

	// Remove leading slash if present
	filePath = strings.TrimPrefix(filePath, "/")

	// Print the requested file path for debugging
	fmt.Printf("Requested file path: %s\n", filePath)

	// Construct the full path to the file in the filesystem
	fullPath := filepath.Join(LaravelStoragePath, filePath)
	
	fmt.Printf("Accessing file at: %s\n", fullPath)
	
	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "File not found",
			"path": fullPath,
		})
		return
	}
	
	// Open the file
	file, err := os.Open(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to open file",
			"details": err.Error(),
			"path": fullPath,
		})
		return
	}
	defer file.Close()
	
	// Get file info for content length
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get file info",
			"details": err.Error(),
		})
		return
	}
	
	// Set content type based on file extension
	extension := strings.ToLower(filepath.Ext(fullPath))
	contentType := "application/octet-stream" // Default
	
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
	
	// Set content disposition for viewing
	filename := filepath.Base(filePath)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))
	
	// Stream the file to the client
	c.Status(http.StatusOK)
	io.Copy(c.Writer, file)
}

// DirectFileDownload serves files directly from the filesystem with download header
func DirectFileDownload(c *gin.Context) {
	// Get the file path from the URL
	filePath := c.Param("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File path is required"})
		return
	}

	// Remove leading slash if present
	filePath = strings.TrimPrefix(filePath, "/")

	// Print the requested file path for debugging
	fmt.Printf("Requested file path for download: %s\n", filePath)

	// Construct the full path to the file in the filesystem
	fullPath := filepath.Join(LaravelStoragePath, filePath)
	
	fmt.Printf("Accessing file at: %s\n", fullPath)
	
	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "File not found",
			"path": fullPath,
		})
		return
	}
	
	// Open the file
	file, err := os.Open(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to open file",
			"details": err.Error(),
			"path": fullPath,
		})
		return
	}
	defer file.Close()
	
	// Get file info for content length
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get file info",
			"details": err.Error(),
		})
		return
	}
	
	// Set content type based on file extension
	extension := strings.ToLower(filepath.Ext(fullPath))
	contentType := "application/octet-stream" // Default
	
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
	
	// Set content disposition for download
	filename := filepath.Base(filePath)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	
	// Stream the file to the client
	c.Status(http.StatusOK)
	io.Copy(c.Writer, file)
}