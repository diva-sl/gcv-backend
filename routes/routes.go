package routes

import (
	"gcv-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Existing Public Routes
		api.POST("/contact", controllers.SubmitContact)
		api.POST("/subscribe", controllers.SubmitSubscribe)

		// Admin Authentication Route
		api.POST("/admin/login", controllers.LoginAdmin)

		// Public Projects Case Study Routes
		api.GET("/projects", controllers.GetProjects)
		api.GET("/projects/:projectId", controllers.GetProjectByID)

		// Protected Admin Actions Group (Requires valid JWT Token)
		admin := api.Group("/admin")
		admin.Use(controllers.RequireAdminAuth())
		{
			admin.GET("/upload/url", controllers.GetPresignedUploadURL)
			admin.POST("/projects", controllers.CreateProject)
			admin.PUT("/projects/:projectId", controllers.UpdateProject)
			admin.DELETE("/projects/:projectId", controllers.DeleteProject)
		}
	}
}