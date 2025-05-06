package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/controllers"
	middleware "github.com/rudychandra/lagi/middlewares"
)

func SetupRouter(r *gin.Engine) {
	// Endpoint login
	r.POST("/login", controllers.Login)

	mahasiswa := r.Group("/bimbingan")
	mahasiswa.Use(middleware.InternalAuthMiddleware()) // pakai token internal
	{
		// Mahasiswa bisa lihat bimbingan mereka (GET)
		mahasiswa.GET("/", controllers.GetBimbingan)
		// Mahasiswa bisa request bimbingan baru (POST)
		mahasiswa.POST("/", controllers.CreateBimbingan)
		// Mahasiswa bisa delete bimbingan mereka sendiri (DELETE)
	}
	
	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Student routes
		api.GET("/student", controllers.GetStudentData)
		
		// Add other protected routes here
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
	
	// Add jadwal routes - THIS WAS MISSING
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
		pengumpulan.GET("/", controllers.GetSubmitanTugas)                   // Get all tasks for the user's group
		pengumpulan.GET("/:id", controllers.GetSubmitanTugasById)            // Get a specific task by ID
		pengumpulan.POST("/:id/upload", controllers.UploadFileTugas)         // Upload a file for a task (first time)
		pengumpulan.PUT("/:id/upload", controllers.UploadFileTugas) // Update/replace uploaded file
	}
	
	// Artefak
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