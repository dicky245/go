package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/routes"
)

func main() {
	// Inisialisasi koneksi ke database dan memuat konfigurasi
	config.Connect()

	// Set up Gin router
	r := gin.Default()
	routes.SetupRouter(r)

	// Port server
	port := "8080"
	fmt.Println("Server berjalan di http://0.0.0.0:" + port)
	err := r.Run("0.0.0.0:" + port)

	if err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
