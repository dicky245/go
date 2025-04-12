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
	pengumuman := r.Group("/pengumuman")
	pengumuman.Use(middleware.InternalAuthMiddleware())
	{
		pengumuman.GET("/", controllers.GetPengumuman)
		pengumuman.GET("/:id", controllers.GetPengumumanByID)
		pengumuman.POST("/", controllers.CreatePengumuman)
		pengumuman.DELETE("/:id", controllers.DeletePengumuman)
	}

	approve := r.Group("/approve")
	approve.Use(middleware.InternalAuthMiddleware())
	{
		approve.GET("/", controllers.GetUpdateBimbingan)
		approve.GET("/:id", controllers.GetUpdateBimbinganByID)
		approve.PUT("/:id", controllers.UpdateRequestBimbingan)
	}
	jadwal := r.Group("/jadwal")
	jadwal.Use(middleware.InternalAuthMiddleware())
	{
		jadwal.GET("/", controllers.GetJadwal)
		jadwal.GET("/:id", controllers.GetJadwalByID)
		jadwal.POST("/", controllers.CreateJadwal)
	}
	pengumpulan := r.Group("/pengumpulan")
	pengumpulan.Use(middleware.InternalAuthMiddleware())
	{
		pengumpulan.GET("/", controllers.GetPengumpulan)
		pengumpulan.GET("/:id", controllers.GetPengumpulanByID)
		pengumpulan.POST("/", controllers.CreatePengumpulan)
		pengumpulan.DELETE("/:id", controllers.DeletePengumpulan)
	}
	// artefak
	artefak := r.Group("/artefak")
	artefak.Use(middleware.InternalAuthMiddleware())
	{
		artefak.GET("/", controllers.GetArtefak)
		artefak.GET("/:id", controllers.GetArtefakByID)
		artefak.POST("/", controllers.CreateArtefak)
		artefak.DELETE("/:id", controllers.DeleteArtefak)
		artefak.GET("/submitan", controllers.GetSubmitMahasiswa)

	}

}
func RoleRoutes(r *gin.Engine) {
	roleGroup := r.Group("/roles")
	{
		roleGroup.POST("/", controllers.CreateRole)
		roleGroup.GET("/", controllers.GetRoles)
		roleGroup.GET("/:id", controllers.GetRoleByID)
		roleGroup.PUT("/:id", controllers.UpdateRole)
		roleGroup.DELETE("/:id", controllers.DeleteRole)
	}
}
func DosenRoleRoutes(r *gin.Engine) {
	roleGroup := r.Group("/dosenroles")
	{
		roleGroup.POST("/", controllers.CreateDosenRoles)
		roleGroup.GET("/", controllers.GetDosenroles)
		roleGroup.PUT("/:id", controllers.UpdateDosenrole)
		roleGroup.DELETE("/:id", controllers.DeleteDosenRole)
		roleGroup.GET("/prodi/:prodi", controllers.GetDosenRolesByProdi)
		roleGroup.GET("/:id", controllers.GetDosenRolesByid)
	}
}
func KelompokRoutes(r *gin.Engine) {
	roleGroup := r.Group("/kelompok")
	{
		roleGroup.POST("/", controllers.CreateKelompok)
		roleGroup.GET("/", controllers.GetKelompok)
		roleGroup.PUT("/:id", controllers.UpdateKelompok)
		roleGroup.DELETE("/:id", controllers.DeleteKelompok)
		roleGroup.GET("/:id", controllers.GetKelompokByID)
	}
}
