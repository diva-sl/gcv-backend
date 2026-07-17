package routes

import (
	"gcv-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/contact", controllers.SubmitContact)
		api.POST("/subscribe", controllers.SubmitSubscribe)
	}
}
