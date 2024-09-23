package routes

import (
	"ShowTimes/pkg/api/handler"
	"ShowTimes/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.RouterGroup, userHandler *handler.UserHandler, otphandler *handler.OtpHandler) {
	engine.GET("/google_callback", userHandler.GoogleCallback)
	engine.GET("/google_login", userHandler.Authv2)

	engine.POST("/signup", userHandler.UserSignUp)
	engine.POST("/login", userHandler.LoginHandler)

	engine.POST("/otplogin", otphandler.SendOTP)
	engine.POST("/verifyotp", otphandler.VerifyOTP)

	engine.Use(middleware.UserAuthMiddleware)
	{

		profile := engine.Group("/profile")
		{
			profile.POST("/addaddress", userHandler.AddAddress)
			profile.GET("/showuserdetails", userHandler.ShowUserDetails)
			profile.GET("/alladdress", userHandler.GetAllAddress)
			profile.PUT("/editprofile", userHandler.EditProfile)
			profile.PATCH("/changepassword",userHandler.ChangePassword)
		}

	}
}
