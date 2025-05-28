package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Configuration for direct file access
const (
	// Update this to the actual path to your uploads directory
	UploadsPath = "D:/Semester 4/PA II/go/uploads/tugas"
)

// MobileFileView serves files directly from the filesystem
func MobileFileView(c *gin.Context) {
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
	fullPath := filepath.Join(UploadsPath, filePath)
	
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
	
	// Set CORS headers to allow direct access from mobile app
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	
	// Set content disposition for viewing
	filename := filepath.Base(filePath)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))
	
	// Stream the file to the client
	c.Status(http.StatusOK)
	io.Copy(c.Writer, file)
}

// MobileFileDownload serves files directly from the filesystem with download header
func MobileFileDownload(c *gin.Context) {
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
	fullPath := filepath.Join(UploadsPath, filePath)
	
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
	
	// Set CORS headers to allow direct access from mobile app
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	
	// Set content disposition for download
	filename := filepath.Base(filePath)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	
	// Stream the file to the client
	c.Status(http.StatusOK)
	io.Copy(c.Writer, file)
}

// MobileTestFileAccess tests if a file can be accessed directly
func MobileTestFileAccess(c *gin.Context) {
	// Test file paths
	testPaths := []string{
		"tugas/4_20250513102545_Doc1.docx",
		"tugas/3_20250513085737_mpdf (1) (9).pdf",
	}
	
	results := make([]map[string]interface{}, 0)
	
	for _, testPath := range testPaths {
		fullPath := filepath.Join(UploadsPath, testPath)
		
		// Check if file exists
		fileInfo, err := os.Stat(fullPath)
		
		result := map[string]interface{}{
			"path": testPath,
			"full_path": fullPath,
			"exists": err == nil,
		}
		
		if err == nil {
			result["size"] = fileInfo.Size()
			result["modified"] = fileInfo.ModTime().Format(time.RFC3339)
			result["is_dir"] = fileInfo.IsDir()
			
			// Generate URLs for testing
			result["view_url"] = fmt.Sprintf("%s/mobile-files/view/%s", c.Request.Host, testPath)
			result["download_url"] = fmt.Sprintf("%s/mobile-files/download/%s", c.Request.Host, testPath)
		} else {
			result["error"] = err.Error()
		}
		
		results = append(results, result)
	}
	
	// Also check the uploads directory itself
	uploadsInfo, err := os.Stat(UploadsPath)
	uploadsResult := map[string]interface{}{
		"path": UploadsPath,
		"exists": err == nil,
	}
	
	if err == nil {
		uploadsResult["is_dir"] = uploadsInfo.IsDir()
		uploadsResult["modified"] = uploadsInfo.ModTime().Format(time.RFC3339)
		
		// List files in the uploads directory
		if uploadsInfo.IsDir() {
			files, err := os.ReadDir(UploadsPath)
			if err == nil {
				fileList := make([]string, 0)
				for _, file := range files {
					fileList = append(fileList, file.Name())
				}
				uploadsResult["files"] = fileList
			} else {
				uploadsResult["list_error"] = err.Error()
			}
		}
	} else {
		uploadsResult["error"] = err.Error()
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Mobile file access test",
		"uploads_path": UploadsPath,
		"test_results": results,
		"uploads_directory": uploadsResult,
		"server_time": time.Now().Format(time.RFC3339),
	})
}

// ListMobileFiles lists files in a directory
func ListMobileFiles(c *gin.Context) {
	// Get the directory path from the query parameter
	dirPath := c.Query("path")
	if dirPath == "" {
		dirPath = "tugas" // Default to tugas directory
	}
	
	// Construct the full path to the directory
	fullPath := filepath.Join(UploadsPath, dirPath)
	
	// Check if directory exists
	dirInfo, err := os.Stat(fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Directory not found",
			"path": fullPath,
			"details": err.Error(),
		})
		return
	}
	
	// Make sure it's a directory
	if !dirInfo.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Path is not a directory",
			"path": fullPath,
		})
		return
	}
	
	// Read directory contents
	files, err := os.ReadDir(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read directory",
			"path": fullPath,
			"details": err.Error(),
		})
		return
	}
	
	// Format the response
	fileList := make([]map[string]interface{}, 0)
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			continue
		}
		
		filePath := filepath.Join(dirPath, file.Name())
		
		fileData := map[string]interface{}{
			"name": file.Name(),
			"path": filePath,
			"is_dir": file.IsDir(),
			"size": fileInfo.Size(),
			"modified": fileInfo.ModTime().Format(time.RFC3339),
		}
		
		// Add URLs for files
		if !file.IsDir() {
			fileData["view_url"] = fmt.Sprintf("/mobile-files/view/%s", filePath)
			fileData["download_url"] = fmt.Sprintf("/mobile-files/download/%s", filePath)
		}
		
		fileList = append(fileList, fileData)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Directory listing",
		"path": dirPath,
		"full_path": fullPath,
		"files": fileList,
	})
}