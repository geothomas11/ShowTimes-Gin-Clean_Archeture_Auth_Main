package routes

import (
	"ShowTimes/pkg/api/handler"
	"ShowTimes/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.RouterGroup, userHandler *handler.UserHandler, otphandler *handler.OtpHandler) {
	engine.POST("/signup", userHandler.UserSignUp)
	engine.POST("/login", userHandler.LoginHandler)

	engine.POST("/otplogin", otphandler.SendOTP)
	engine.POST("/verifyotp", otphandler.VerifyOTP)

	engine.Use(middleware.UserAuthMiddleware)
	{

	}
}
