package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/controllers"
	middleware "github.com/rudychandra/lagi/middlewares"
)

func SetupRouter(r *gin.Engine) {
	// --- Auth ---
	r.POST("/login", controllers.Login)

	// --- Bimbingan (Mahasiswa) ---
	mahasiswa := r.Group("/bimbingan")
	mahasiswa.Use(middleware.InternalAuthMiddleware())
	{
		mahasiswa.GET("/", controllers.GetBimbingan)
		mahasiswa.POST("/", controllers.CreateBimbingan)
	}

	// --- Ruangan (Tanpa Auth, untuk dropdown) ---
	r.GET("/ruangans", controllers.GetRuangans)

	// --- API / Protected Routes (Authenticated user) ---
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/student", controllers.GetStudentData)
		// Tambahkan API lainnya di sini
	}

	// --- Pengumuman (Mahasiswa + Dosen) ---
	pengumuman := r.Group("/pengumuman")
	pengumuman.Use(middleware.InternalAuthMiddleware())
	{
		pengumuman.GET("/", controllers.GetPengumuman)
		pengumuman.GET("/:id", controllers.GetPengumumanByID)
		pengumuman.POST("/", controllers.CreatePengumuman)
		pengumuman.DELETE("/:id", controllers.DeletePengumuman)
	}

	// --- Approval Request Bimbingan (Dosen) ---
	approve := r.Group("/approve")
	approve.Use(middleware.InternalAuthMiddleware())
	{
		approve.GET("/", controllers.GetUpdateBimbingan)
		approve.GET("/:id", controllers.GetUpdateBimbinganByID)
		approve.PUT("/:id", controllers.UpdateRequestBimbingan)
	}

	// --- Jadwal (Mahasiswa + Dosen) ---
	jadwal := r.Group("/jadwal")
	jadwal.Use(middleware.InternalAuthMiddleware())
	{
		jadwal.GET("/", controllers.GetJadwal)
		jadwal.GET("/:id", controllers.GetJadwalByID)
		jadwal.POST("/", controllers.CreateJadwal)
		jadwal.GET("/ruangan", controllers.GetRuangan) // Opsional/alternatif endpoint
	}

	// --- Pengumpulan Tugas ---
	pengumpulan := r.Group("/pengumpulan")
	pengumpulan.Use(middleware.InternalAuthMiddleware())
	{
		pengumpulan.GET("/", controllers.GetAllTugas)
		pengumpulan.GET("/:id", controllers.GetTugasById)
		pengumpulan.POST("/:id/upload", controllers.CreateTugas)
		pengumpulan.PUT("/:id/upload", controllers.UpdateTugas)
	}

	// --- Tugas Management ---
	TugasRoutes(r)
	// --- Role Management ---
	RoleRoutes(r)
	DosenRoleRoutes(r)
}

func TugasRoutes(r *gin.Engine) {
	tugasGroup := r.Group("/tugas")
	tugasGroup.Use(middleware.InternalAuthMiddleware()) // Using the same auth middleware as other routes
	{
		tugasGroup.GET("/", controllers.GetAllTugas)
		tugasGroup.GET("/:id", controllers.GetTugasById)
		tugasGroup.POST("/", controllers.CreateTugas)
		tugasGroup.PUT("/:id", controllers.UpdateTugas)
		tugasGroup.DELETE("/:id", controllers.DeleteTugas)
		tugasGroup.GET("/kategori/:kategori", controllers.GetTugasByKategori)
		tugasGroup.GET("/status/:status", controllers.GetTugasByStatus)
		tugasGroup.GET("/prodi/:prodi_id", controllers.GetTugasByProdi)
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