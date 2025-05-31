package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
)

// CreateUser saves or updates device token (legacy endpoint untuk manual create)
func CreateUser(c *gin.Context) {
	// Get database connection
	db, err := config.GetDB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database connection not available"})
		return
	}

	var newUser model.Device_Token

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data tidak valid: " + err.Error(),
		})
		return
	}

	var exists model.Device_Token
	result := db.Where("user_id = ?", newUser.UserID).First(&exists)

	if result.RowsAffected > 0 {
		exists.TokenDevice = newUser.TokenDevice
		if err := db.Save(&exists).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal memperbaharui token: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Token sudah diperbaharui",
			"data":    exists,
		})
		return
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat user baru: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User berhasil ditambahkan",
		"data":    newUser,
	})
}

// GetDeviceToken retrieves device token for the authenticated user
func GetDeviceToken(c *gin.Context) {
	// Get user ID from token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get database connection
	db, err := config.GetDB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database connection not available"})
		return
	}

	var deviceToken model.Device_Token
	if err := db.Where("user_id = ?", userID).First(&deviceToken).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Device token not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   deviceToken,
	})
}

// SaveDeviceToken saves or updates device token for the authenticated user
func SaveDeviceToken(c *gin.Context) {
	// Get user ID from token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get database connection
	db, err := config.GetDB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database connection not available"})
		return
	}

	var request struct {
		TokenDevice string `json:"token_device" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token device is required"})
		return
	}

	// Convert userID to int
	userIDInt, ok := userID.(int)
	if !ok {
		// Try to convert from uint or other types
		if userIDUint, ok := userID.(uint); ok {
			userIDInt = int(userIDUint)
		} else if userIDFloat, ok := userID.(float64); ok {
			userIDInt = int(userIDFloat)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
			return
		}
	}

	// Check if device token already exists for this user
	var deviceToken model.Device_Token
	result := db.Where("user_id = ?", userIDInt).First(&deviceToken)

	if result.Error != nil {
		// Create new device token
		newDeviceToken := model.Device_Token{
			UserID:      userIDInt,
			TokenDevice: request.TokenDevice,
		}

		if err := db.Create(&newDeviceToken).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save device token"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "Device token saved successfully",
			"data":    newDeviceToken,
		})
	} else {
		// Update existing device token
		deviceToken.TokenDevice = request.TokenDevice

		if err := db.Save(&deviceToken).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Device token updated successfully",
			"data":    deviceToken,
		})
	}
}

// DeleteDeviceToken removes device token for the authenticated user
func DeleteDeviceToken(c *gin.Context) {
	// Get user ID from token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get database connection
	db, err := config.GetDB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database connection not available"})
		return
	}

	if err := db.Where("user_id = ?", userID).Delete(&model.Device_Token{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete device token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Device token deleted successfully",
	})
}
