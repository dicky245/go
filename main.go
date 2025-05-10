	package main

	import (
		"fmt"
		"log"

		"github.com/gin-gonic/gin"
		"github.com/rudychandra/lagi/config"
		"github.com/rudychandra/lagi/routes"
	)

	func main() {
		config.Connect()
	
		r := gin.Default()
		routes.SetupRouter(r) // Sudah cukup, jangan panggil RoleRoutes lagi
		
	
		port := "8080"
		fmt.Println("Server berjalan di http://0.0.0.0:" + port)
		err := r.Run("0.0.0.0:" + port)
	
		if err != nil {
			log.Fatal("Gagal menjalankan server:", err)
		}
	}
	
	
