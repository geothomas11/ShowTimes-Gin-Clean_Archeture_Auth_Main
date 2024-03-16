package routes

import (
	"ShowTimes/pkg/api/handler"
	"ShowTimes/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(engine *gin.RouterGroup, adminHandler *handler.AdminHandler) {
	engine.POST("/adminlogin", adminHandler.LoginHandler)

	engine.Use(middleware.AdminAuthMiddleware)
	{
		userManagement := engine.Group("/users")
		{
			userManagement.PUT("/block", adminHandler.BlockUser)
			userManagement.PUT("/unblock", adminHandler.UnBlockUser)
			userManagement.GET("", adminHandler.GetUsers)

		}
	}
}
