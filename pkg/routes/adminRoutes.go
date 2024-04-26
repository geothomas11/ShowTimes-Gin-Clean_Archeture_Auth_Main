package routes

import (
	"ShowTimes/pkg/api/handler"
	"ShowTimes/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(engine *gin.RouterGroup, adminHandler *handler.AdminHandler, CategoryHandler *handler.CategoryHandler, inventoryHandler *handler.InventoryHandler) {
	engine.POST("/adminlogin", adminHandler.LoginHandler)

	engine.Use(middleware.AdminAuthMiddleware)
	{
		userManagement := engine.Group("/users")
		{
			userManagement.PUT("/block", adminHandler.BlockUser)
			userManagement.PUT("/unblock", adminHandler.UnBlockUser)
			userManagement.GET("", adminHandler.GetUsers)

		}
		categorymanagement := engine.Group("/category")

		{
			categorymanagement.POST("/addcategory", CategoryHandler.AddCategory)
			categorymanagement.GET("/getcategory", CategoryHandler.GetCategory)
			categorymanagement.PUT("/updatecategory", CategoryHandler.UpdateCategory)
			categorymanagement.DELETE("/deletecategory", CategoryHandler.DeleteCategory)
		}
		inventorymanagement := engine.Group("/inventory")

		{
			inventorymanagement.POST("/addinventory", inventoryHandler.AddInventory)
			inventorymanagement.GET("/listinventory", inventoryHandler.ListProducts)
			inventorymanagement.PUT("/editinventory", inventoryHandler.EditInventory)
			inventorymanagement.DELETE("/deleteinventory", inventoryHandler.DeleteInventory)
			inventorymanagement.PATCH("/updateinventory", inventoryHandler.UpdateInventory)
		}
	}
}
