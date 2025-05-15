package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// DebugFileAccess helps debug file access issues
func DebugFileAccess(c *gin.Context) {
	// Get the file path from the URL
	filePath := c.Query("path")
	if filePath == "" {
		filePath = "tugas_files/xTP3hxTxPwmeU6O9wJuxCLbIJbFAfapwstSljFmV.pdf" // Default test file
	}

	// Test different URL constructions
	urls := []string{
		fmt.Sprintf("%s/%s%s", LaravelBaseURL, StorageBasePath, filePath),
		fmt.Sprintf("%s/%s", LaravelBaseURL, filePath),
		fmt.Sprintf("%s/storage/%s", LaravelBaseURL, filePath),
		fmt.Sprintf("%s/public/storage/%s", LaravelBaseURL, filePath),
	}

	results := make([]map[string]interface{}, 0)

	// Test each URL
	for _, url := range urls {
		client := &http.Client{}
		req, err := http.NewRequest("HEAD", url, nil)
		if err != nil {
			results = append(results, map[string]interface{}{
				"url":     url,
				"success": false,
				"error":   "Failed to create request: " + err.Error(),
			})
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			results = append(results, map[string]interface{}{
				"url":     url,
				"success": false,
				"error":   "Failed to connect: " + err.Error(),
			})
			continue
		}
		defer resp.Body.Close()

		results = append(results, map[string]interface{}{
			"url":         url,
			"success":     resp.StatusCode == http.StatusOK,
			"status_code": resp.StatusCode,
			"headers":     resp.Header,
		})
	}

	// Check local storage paths
	localPaths := []string{
		"storage/app/public/" + filePath,
		"public/storage/" + filePath,
		"storage/" + filePath,
	}

	localResults := make([]map[string]interface{}, 0)
	for _, path := range localPaths {
		_, err := os.Stat(path)
		var errorStr string
		if err != nil {
			errorStr = err.Error()
		}
		
		absPath, absErr := filepath.Abs(path)
		if absErr != nil {
			absPath = "Error getting absolute path: " + absErr.Error()
		}
		
		localResults = append(localResults, map[string]interface{}{
			"path":    path,
			"exists":  err == nil,
			"error":   errorStr,
			"absPath": absPath,
		})
	}

	// Get current directory with proper error handling
	currentDir, currentDirErr := filepath.Abs(".")
	if currentDirErr != nil {
		currentDir = "Error getting current directory: " + currentDirErr.Error()
	}

	// Get working directory with proper error handling
	workingDir, workingDirErr := os.Getwd()
	if workingDirErr != nil {
		workingDir = "Error getting working directory: " + workingDirErr.Error()
	}

	c.JSON(http.StatusOK, gin.H{
		"file_path":     filePath,
		"url_tests":     results,
		"local_tests":   localResults,
		"laravel_base":  LaravelBaseURL,
		"storage_base":  StorageBasePath,
		"current_dir":   currentDir,
		"working_dir":   workingDir,
		"env_variables": os.Environ(),
	})
}