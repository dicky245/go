package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/controllers"
	middleware "github.com/rudychandra/lagi/middlewares"
)

func SetupRouter(r *gin.Engine) {
	// Endpoint login
	r.POST("/login", controllers.Login)

	// Group untuk endpoint internal mahasiswa - menggunakan token internal
	mahasiswa := r.Group("/bimbingan")
	mahasiswa.Use(middleware.InternalAuthMiddleware()) // pakai token internal
	{
		mahasiswa.GET("/", controllers.GetBimbingan)
		mahasiswa.POST("/", controllers.CreateBimbingan)
		mahasiswa.DELETE("/:id", controllers.DeleteBimbingan)
	}
	// dosen := r.Group("/dosen")
	// dosen.Use(middleware.InternalAuthMiddleware())
	// {
	//     dosen.GET("/jadwal", controllers.GetJadwalDosen)
	// }
}
