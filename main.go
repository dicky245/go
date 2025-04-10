package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/routes"
)

func main() {
	
	config.Connect() // Hubungkan database
	

	r := gin.Default()
	routes.SetupRouter(r)
	routes.RoleRoutes(r)
	routes.DosenRoleRoutes(r)
	routes.KelompokRoutes(r)
	

	port := "8080"
	fmt.Println("Server berjalan di http://localhost:" + port)
	err := r.Run(":" + port)

	if err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
