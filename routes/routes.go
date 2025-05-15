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
	// Use the dedicated setup function for pengumpulan routes
	SetupPengumpulanRoutes(r)

	// --- Tugas Management ---
	TugasRoutes(r)

	// --- Role Management ---
	RoleRoutes(r)
	DosenRoleRoutes(r)
	SetupFileRoutes(r) // Add this line to include file routes
}

// SetupPengumpulanRoutes configures all routes related to assignment submissions
func SetupPengumpulanRoutes(r *gin.Engine) {
	pengumpulan := r.Group("/pengumpulan")
	pengumpulan.Use(middleware.InternalAuthMiddleware())
	{
		// Get all tugas for mahasiswa
		pengumpulan.GET("/", controllers.GetSubmitanTugas)
		
		// Get specific tugas by ID
		pengumpulan.GET("/:id", controllers.GetSubmitanTugasById)
		
		// Create a new submission for an assignment
		pengumpulan.POST("/:id/upload", controllers.UpdateUploadFileTugas)
		
		// Update an existing submission - Add this missing route
		pengumpulan.PUT("/:id/upload", controllers.UpdateUploadFileTugas)
	}
}

func TugasRoutes(r *gin.Engine) {
	tugasGroup := r.Group("/tugas")
	tugasGroup.Use(middleware.InternalAuthMiddleware()) // Using the same auth middleware as other routes
	{
		tugasGroup.GET("/", controllers.GetSubmitanTugas)             // Get all tugas for mahasiswa
		tugasGroup.GET("/:id", controllers.GetSubmitanTugasById)         // Get specific tugas by ID
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
// Add this to your SetupFileRoutes function
func SetupFileRoutes(r *gin.Engine) {
	// Web file routes (Laravel storage proxy)
	webFiles := r.Group("/files")
	{
		// View file (proxies to Laravel storage)
		webFiles.GET("/view/*path", controllers.ProxyFile)
		
		// Download file (with attachment header)
		webFiles.GET("/download/*path", controllers.DownloadFile)
		
		// Debug endpoint
		webFiles.GET("/debug", controllers.DebugFileAccess)
		webFiles.GET("/test", controllers.TestLaravelStorage)
	}

	// Mobile file routes (direct file access)
	mobileFiles := r.Group("/mobile-files")
	{
		// View file directly from filesystem
		mobileFiles.GET("/view/*path", controllers.MobileFileView)
		
		// Download file directly from filesystem
		mobileFiles.GET("/download/*path", controllers.MobileFileDownload)
		
		// Test endpoint
		mobileFiles.GET("/test", controllers.MobileTestFileAccess)
		
		// List files in directory
		mobileFiles.GET("/list", controllers.ListMobileFiles)
	}

	// Protected file routes (auth required)
	protectedFiles := r.Group("/secure-files")
	protectedFiles.Use(middleware.InternalAuthMiddleware())
	{
		// These could be added if you need additional file operations
		// that require authentication
		protectedFiles.GET("/view/*path", controllers.ProxyFile)
		protectedFiles.GET("/download/*path", controllers.DownloadFile)
	}
}