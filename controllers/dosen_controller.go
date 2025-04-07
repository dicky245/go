package controllers

import (
	"net/http"
	"strconv"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"

	"github.com/gin-gonic/gin"
)
	
type DosenRoleRequest struct {
	UserID  uint   `json:"user_id"`
	Prodi   string   `json:"prodi"`
	Tingkat uint     `json:"tingkat"`
	RoleIDs []uint `json:"role_ids"`
}

func CreateDosenRoles(c *gin.Context) {
	var req DosenRoleRequest

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var createdRoles []model.DosenRole // Menyimpan data yang berhasil dibuat

	for _, roleID := range req.RoleIDs {
		dosenRole := model.DosenRole{
			UserID: req.UserID,
			Prodi: req.Prodi,
			Tingkat: req.Tingkat,
			RoleID: roleID,
		}

		// Cek apakah kombinasi user_id + role_id sudah ada
		var exists model.DosenRole
		if err := config.DB.
			Where("user_id = ? AND prodi = ? AND tingkat = ? AND role_id = ?", req.UserID, req.Prodi,req.Tingkat, roleID).
			First(&exists).Error; err == nil {
			// Skip jika sudah ada
			continue
		}

		// Simpan ke database
		if err := config.DB.Create(&dosenRole).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		createdRoles = append(createdRoles, dosenRole)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Roles berhasil ditambahkan",
		"data":    createdRoles,
	})
}

// Get All Roles
func GetDosenroles(c *gin.Context) {
	var dosen_roles []model.DosenRole
	config.DB.Find(&dosen_roles)
	c.JSON(http.StatusOK, dosen_roles)
}

func UpdateDosenrole(c *gin.Context) {
	var dosenRole model.DosenRole

	// **Konversi id dari string ke uint**
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// **Cek apakah data ada**
	if err := config.DB.First(&dosenRole, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	// **Bind data JSON ke struct sementara**
	var req model.DosenRole
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// **Update field yang diperlukan**
	dosenRole.UserID = req.UserID
	dosenRole.RoleID = req.RoleID
	dosenRole.Prodi = req.Prodi
	dosenRole.Tingkat =req.Tingkat

	// **Simpan perubahan**
	config.DB.Save(&dosenRole)
	c.JSON(http.StatusOK, gin.H{"message": "Role updated", "data": dosenRole})
}

// Get DosenRoles by Prodi
func GetDosenRolesByProdi(c *gin.Context) {
	prodi := c.Param("prodi")

	var dosenRoles []model.DosenRole
	if err := config.DB.Where("prodi = ?", prodi).Find(&dosenRoles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dosenRoles)
}
func GetDosenRolesByid(c *gin.Context) {
	id := c.Param("id")

	var dosenRoles []model.DosenRole
	if err := config.DB.Where("id = ?", id).Find(&dosenRoles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dosenRoles)
}

// 🔹 **DELETE DOSEN ROLE (FIX BUG!)**
func DeleteDosenRole(c *gin.Context) {
	var dosenRole model.DosenRole

	// **Konversi id dari string ke uint**
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// **Cek apakah data ada**
	if err := config.DB.First(&dosenRole, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	// **Hapus dari database**
	config.DB.Delete(&dosenRole)
	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}