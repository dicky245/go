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
		mahasiswa.GET("/", controllers.GetBimbingan)       // Mahasiswa bisa lihat bimbingan mereka
		mahasiswa.POST("/", controllers.CreateBimbingan)   // Mahasiswa bisa request bimbingan
		mahasiswa.DELETE("/:id", controllers.DeleteBimbingan) // Mahasiswa bisa delete bimbingan
	}
	
	pengumuman := r.Group("/pengumuman")
	pengumuman.Use(middleware.InternalAuthMiddleware())
	{
		pengumuman.GET("/", controllers.GetPengumuman)   // Mahasiswa bisa lihat pengumuman
		pengumuman.GET("/:id", controllers.GetPengumumanByID) // Mahasiswa bisa lihat detail pengumuman
		pengumuman.POST("/", controllers.CreatePengumuman)  // Dosen buat pengumuman
		pengumuman.DELETE("/:id", controllers.DeletePengumuman) // Dosen delete pengumuman
	}
		
	approve := r.Group("/approve")
	approve.Use(middleware.InternalAuthMiddleware())
	{
		approve.GET("/", controllers.GetUpdateBimbingan) // Dosen bisa approve bimbingan
		approve.GET("/:id", controllers.GetUpdateBimbinganByID) // Dosen lihat detail
		approve.PUT("/:id", controllers.UpdateRequestBimbingan) // Dosen update request bimbingan
	}
	
	jadwal := r.Group("/jadwal")
	jadwal.Use(middleware.InternalAuthMiddleware())
	{
		jadwal.GET("/", controllers.GetJadwal)       // Mahasiswa lihat jadwal
		jadwal.GET("/:id", controllers.GetJadwalByID) // Mahasiswa lihat jadwal spesifik
		jadwal.POST("/", controllers.CreateJadwal)    // Dosen buat jadwal
	}
	
	pengumpulan := r.Group("/pengumpulan")
	pengumpulan.Use(middleware.InternalAuthMiddleware())
	{
		pengumpulan.GET("/", controllers.GetPengumpulan)
		pengumpulan.GET("/:id", controllers.GetPengumpulanByID)
		pengumpulan.POST("/", controllers.CreatePengumpulan)
		// pengumpulan.PUT("/:id", controllers.UpdatePengumpulan)
		pengumpulan.DELETE("/:id", controllers.DeletePengumpulan)
		pengumpulan.POST("/:id/upload", controllers.UploadFilePengumpulan) // Upload file endpoint
	}
	
	

	// Artefak
	artefak := r.Group("/artefak")
	artefak.Use(middleware.InternalAuthMiddleware())
	{
		artefak.GET("/", controllers.GetArtefak)      // Mahasiswa dan Dosen lihat artefak
		artefak.GET("/:id", controllers.GetArtefakByID) // Mahasiswa dan Dosen lihat artefak detail
		artefak.POST("/", controllers.CreateArtefak)   // Dosen buat artefak
		artefak.DELETE("/:id", controllers.DeleteArtefak) // Dosen hapus artefak
		artefak.GET("/submitan", controllers.GetSubmitMahasiswa) // Mahasiswa lihat submitan mereka
	}
}

// Rute untuk role
func RoleRoutes(r *gin.Engine) {
	roleGroup := r.Group("/roles")
	{
		roleGroup.POST("/", controllers.CreateRole)  // Admin buat role baru
		roleGroup.GET("/", controllers.GetRoles)    // Admin lihat semua role
		roleGroup.GET("/:id", controllers.GetRoleByID) // Admin lihat role tertentu
		roleGroup.PUT("/:id", controllers.UpdateRole)  // Admin update role
		roleGroup.DELETE("/:id", controllers.DeleteRole) // Admin hapus role
	}
}

// Rute untuk dosen roles
func DosenRoleRoutes(r *gin.Engine) {
	roleGroup := r.Group("/dosenroles")
	{
		roleGroup.POST("/", controllers.CreateDosenRoles)  // Admin buat dosen role
		roleGroup.GET("/", controllers.GetDosenroles)    // Admin lihat dosen role
		roleGroup.PUT("/:id", controllers.UpdateDosenrole) // Admin update dosen role
		roleGroup.DELETE("/:id", controllers.DeleteDosenRole) // Admin hapus dosen role
		roleGroup.GET("/prodi/:prodi", controllers.GetDosenRolesByProdi) // Admin lihat dosen berdasarkan prodi
		roleGroup.GET("/:id", controllers.GetDosenRolesByid) // Admin lihat dosen role berdasarkan id
	}
}

// Rute untuk kelompok
func KelompokRoutes(r *gin.Engine) {
	roleGroup := r.Group("/kelompok")
	{
		roleGroup.POST("/", controllers.CreateKelompok)  // Admin buat kelompok
		roleGroup.GET("/", controllers.GetKelompok)    // Admin lihat semua kelompok
		roleGroup.PUT("/:id", controllers.UpdateKelompok) // Admin update kelompok
		roleGroup.DELETE("/:id", controllers.DeleteKelompok) // Admin hapus kelompok
		roleGroup.GET("/:id", controllers.GetKelompokByID) // Admin lihat kelompok berdasarkan id
	}
}
