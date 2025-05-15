package controllers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Configuration for Laravel storage
const (
	LaravelBaseURL  = "https://vokasitera-main-ltziwn.laravel.cloud" // Your Laravel app URL
	StorageBasePath = "storage/"              // Laravel's public storage path
)

// ProxyFile handles proxying file requests to Laravel storage
func ProxyFile(c *gin.Context) {
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

	// Construct the full URL to the file in Laravel storage
	fileURL := fmt.Sprintf("%s/%s%s", LaravelBaseURL, StorageBasePath, filePath)
	
	fmt.Printf("Proxying file request to: %s\n", fileURL)
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// Create request
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request", "details": err.Error()})
		return
	}

	// Add headers for debugging
	req.Header.Add("User-Agent", "GoBackendProxy/1.0")
	
	// Send the request
	fmt.Printf("Sending request to: %s\n", fileURL)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching file: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch file", 
			"details": err.Error(),
			"url": fileURL,
		})
		return
	}
	defer resp.Body.Close()

	// Print response status for debugging
	fmt.Printf("Response status: %d\n", resp.StatusCode)
	fmt.Printf("Response headers: %v\n", resp.Header)

	// Check if file was found
	if resp.StatusCode != http.StatusOK {
		// Try to read response body for more details
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Printf("Error response body: %s\n", bodyString)
		
		c.JSON(resp.StatusCode, gin.H{
			"error": "File not found or access denied", 
			"status_code": resp.StatusCode,
			"url": fileURL,
			"response": bodyString,
		})
		return
	}

	// Get content type and set it in the response
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		c.Header("Content-Type", contentType)
	} else {
		// Default to PDF if content type is not provided
		c.Header("Content-Type", "application/pdf")
	}

	// Set content disposition for viewing
	filename := filepath.Base(filePath)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))

	// Stream the file to the client
	c.Status(http.StatusOK)
	bytesWritten, err := io.Copy(c.Writer, resp.Body)
	if err != nil {
		fmt.Printf("Error streaming file to client: %v\n", err)
	} else {
		fmt.Printf("Successfully streamed %d bytes to client\n", bytesWritten)
	}
}

// DownloadFile handles file downloads with attachment disposition
func DownloadFile(c *gin.Context) {
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

	// Construct the full URL to the file in Laravel storage
	fileURL := fmt.Sprintf("%s/%s%s", LaravelBaseURL, StorageBasePath, filePath)
	
	fmt.Printf("Downloading file from: %s\n", fileURL)
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// Create request
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request", "details": err.Error()})
		return
	}

	// Add headers for debugging
	req.Header.Add("User-Agent", "GoBackendProxy/1.0")
	
	// Send the request
	fmt.Printf("Sending request to: %s\n", fileURL)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching file: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch file", 
			"details": err.Error(),
			"url": fileURL,
		})
		return
	}
	defer resp.Body.Close()

	// Print response status for debugging
	fmt.Printf("Response status: %d\n", resp.StatusCode)
	fmt.Printf("Response headers: %v\n", resp.Header)

	// Check if file was found
	if resp.StatusCode != http.StatusOK {
		// Try to read response body for more details
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Printf("Error response body: %s\n", bodyString)
		
		c.JSON(resp.StatusCode, gin.H{
			"error": "File not found or access denied", 
			"status_code": resp.StatusCode,
			"url": fileURL,
			"response": bodyString,
		})
		return
	}

	// Get content type and set it in the response
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		c.Header("Content-Type", contentType)
	} else {
		// Default to PDF if content type is not provided
		c.Header("Content-Type", "application/pdf")
	}

	// Set content disposition for download
	filename := filepath.Base(filePath)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Stream the file to the client
	c.Status(http.StatusOK)
	bytesWritten, err := io.Copy(c.Writer, resp.Body)
	if err != nil {
		fmt.Printf("Error streaming file to client: %v\n", err)
	} else {
		fmt.Printf("Successfully streamed %d bytes to client\n", bytesWritten)
	}
}

// TestLaravelStorage tests the connection to Laravel storage
func TestLaravelStorage(c *gin.Context) {
	// Test different URLs to find which one works
	testURLs := []string{
		fmt.Sprintf("%s/%stugas_files/Dtg7lUEjFMb0B5ur7jCQKrWPA0FEMAfMXN5txfS9.pdf", LaravelBaseURL, StorageBasePath),
		fmt.Sprintf("%s/storage/tugas_files/Dtg7lUEjFMb0B5ur7jCQKrWPA0FEMAfMXN5txfS9.pdf", LaravelBaseURL),
		fmt.Sprintf("%s/public/storage/tugas_files/Dtg7lUEjFMb0B5ur7jCQKrWPA0FEMAfMXN5txfS9.pdf", LaravelBaseURL),
		// Try with IP address instead of localhost
		"https://vokasitera-main-ltziwn.laravel.cloud/storage/tugas_files/Dtg7lUEjFMb0B5ur7jCQKrWPA0FEMAfMXN5txfS9.pdf",
		// Try with your machine's actual IP address
		// Replace 192.168.1.100 with your actual IP address
		// "http://192.168.1.100:8000/storage/tugas_files/xTP3hxTxPwmeU6O9wJuxCLbIJbFAfapwstSljFmV.pdf",
	}
	
	results := []map[string]interface{}{}
	
	for _, url := range testURLs {
		fmt.Printf("Testing URL: %s\n", url)
		
		// Create HTTP client with timeout
		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		
		// Create request
		req, err := http.NewRequest("HEAD", url, nil)
		if err != nil {
			results = append(results, map[string]interface{}{
				"url": url,
				"success": false,
				"error": fmt.Sprintf("Failed to create request: %v", err),
			})
			continue
		}
		
		// Add headers for debugging
		req.Header.Add("User-Agent", "GoBackendProxy/1.0")
		
		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			results = append(results, map[string]interface{}{
				"url": url,
				"success": false,
				"error": fmt.Sprintf("Failed to connect: %v", err),
			})
			continue
		}
		defer resp.Body.Close()
		
		results = append(results, map[string]interface{}{
			"url": url,
			"success": resp.StatusCode == http.StatusOK,
			"status_code": resp.StatusCode,
			"headers": resp.Header,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"laravel_base": LaravelBaseURL,
		"storage_base": StorageBasePath,
	})
}

// GetLaravelStorageFile retrieves a file from Laravel storage and returns the local path
func GetLaravelStorageFile(filePath string) (string, error) {
	// Ensure the cache directory exists
	cacheDir := "cache/files"
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Construct the full URL to the file in Laravel storage
	fileURL := fmt.Sprintf("%s/%s%s", LaravelBaseURL, StorageBasePath, filePath)
	
	// Create HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file: %w", err)
	}
	defer resp.Body.Close()

	// Check if file was found
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("file not found or access denied: %d", resp.StatusCode)
	}

	// Create the cache file
	filename := filepath.Base(filePath)
	cachePath := filepath.Join(cacheDir, filename)
	file, err := os.Create(cachePath)
	if err != nil {
		return "", fmt.Errorf("failed to create cache file: %w", err)
	}
	defer file.Close()

	// Copy the file content
	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", fmt.Errorf("failed to write cache file: %w", err)
	}

	return cachePath, nil
}