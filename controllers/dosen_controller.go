package controllers

import (
	"net/http"
	"strconv"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"

	"github.com/gin-gonic/gin"
)
	
type DosenRoleRequest struct {
	UserID     uint     `json:"user_id"`
	Prodi      string   `json:"prodi"`
	NamaDosen  string   `json:"nama_dosen"`
	Tingkat    uint     `json:"tingkat"`
	RoleID    uint   `json:"role_id"`
	NamaRole  string `json:"nama_role"` // Sesuai jumlah role_ids
}

func CreateDosenRoles(c *gin.Context) {
	var req DosenRoleRequest

	// Ambil input JSON dari request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah data role sudah ada sebelumnya
	var exists model.DosenRole
	if err := config.DB.
		Where("user_id = ? AND prodi = ? AND tingkat = ? AND role_id = ?", req.UserID, req.Prodi, req.Tingkat, req.RoleID).
		First(&exists).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Role sudah ada"})
		return
	}

	// Simpan ke database
	dosenRole := model.DosenRole{
		UserID:    req.UserID,
		NamaDosen: req.NamaDosen,
		RoleID:    req.RoleID,
		NamaRole:  req.NamaRole,
		Prodi:     req.Prodi,
		Tingkat:   req.Tingkat,
	}

	if err := config.DB.Create(&dosenRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Role berhasil ditambahkan",
		"data":    dosenRole,
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
	dosenRole.NamaDosen = req.NamaDosen
	dosenRole.NamaRole = req.NamaRole
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