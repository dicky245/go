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
