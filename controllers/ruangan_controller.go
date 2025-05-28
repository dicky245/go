package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
)

func GetRuangans(c *gin.Context) {
	var ruangans []model.Ruangan

	if err := config.DB.Find(&ruangans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Gagal mengambil data ruangan",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data ruangan berhasil diambil",
		"status":  "success",
		"data":    ruangans,
	})
}
