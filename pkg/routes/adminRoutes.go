package routes

import (
	"ShowTimes/pkg/api/handler"
	"ShowTimes/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(engine *gin.RouterGroup, adminHandler *handler.AdminHandler, CategoryHandler *handler.CategoryHandler, InventoryHandler *handler.ProductHandler) {
	engine.POST("/adminlogin", adminHandler.LoginHandler)

	engine.Use(middleware.AdminAuthMiddleware)
	{
		userManagement := engine.Group("/users")
		{
			userManagement.PUT("/block", adminHandler.BlockUser)
			userManagement.PUT("/unblock", adminHandler.UnBlockUser)
			userManagement.GET("/getusers", adminHandler.GetUsers)

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
			inventorymanagement.POST("/addproducts", InventoryHandler.AddProducts)
			inventorymanagement.GET("/listproducts", InventoryHandler.ListProducts)
			inventorymanagement.PUT("/editproducts", InventoryHandler.EditProducts)
			inventorymanagement.DELETE("/deleteproducts", InventoryHandler.DeleteProducts)
			inventorymanagement.PATCH("/updateproducts", InventoryHandler.UpdateProducts)
		}
	}
}
