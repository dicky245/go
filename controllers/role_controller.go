package controllers

import (
    "github.com/gin-gonic/gin"
   "github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
    "net/http"
    "strconv"
)

// Create Role
func CreateRole(c *gin.Context) {
    var role model.Role
    if err := c.ShouldBindJSON(&role); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := config.DB.Create(&role).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, role)
}

// Get All Roles
func GetRoles(c *gin.Context) {
    var roles []model.Role
    config.DB.Find(&roles)
    c.JSON(http.StatusOK, roles)
}

// Get Role By ID
func GetRoleByID(c *gin.Context) {
    var role model.Role

    // Konversi id dari string ke uint
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    if err := config.DB.First(&role, uint(id)).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
        return
    }
    c.JSON(http.StatusOK, role)
}

// Update Role
func UpdateRole(c *gin.Context) {
    var role model.Role

    // Konversi id dari string ke uint
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    if err := config.DB.First(&role, uint(id)).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
        return
    }

    if err := c.ShouldBindJSON(&role); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    role.ID = uint(id) // Pastikan ID tidak berubah
    config.DB.Save(&role)
    c.JSON(http.StatusOK, role)
}

// Delete Role
func DeleteRole(c *gin.Context) {
    var role model.Role

    // Konversi id dari string ke uint
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    if err := config.DB.First(&role, uint(id)).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
        return
    }

    config.DB.Delete(&role)
    c.JSON(http.StatusOK, gin.H{"message": "Role deleted"})
}