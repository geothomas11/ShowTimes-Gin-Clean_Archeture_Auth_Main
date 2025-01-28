package http

import (
	"ShowTimes/pkg/api/handler"
	"ShowTimes/pkg/routes"
	"log"

	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	Engine *gin.Engine
}

func NewServerHTTP(adminHandler *handler.AdminHandler, userHandler *handler.UserHandler, otpHandler *handler.OtpHandler, Categoryhandler *handler.CategoryHandler, producthandler *handler.ProductHandler, cartHandler *handler.CartHandler, orderHandler *handler.OrderHandler, paymentHandler *handler.PaymentHandler) *ServerHTTP {
	engine := gin.New()

	engine.Use(gin.Logger())

	engine.GET("/validate_token", adminHandler.ValidateRefreshTokenAndCreateNewAccess)

	routes.UserRoutes(engine.Group("/user"), userHandler, otpHandler, cartHandler, orderHandler)
	routes.AdminRoutes(engine.Group("/admin"), adminHandler, Categoryhandler, producthandler, paymentHandler, orderHandler)

	return &ServerHTTP{Engine: engine}
}
func (sh *ServerHTTP) Start() {
	err := sh.Engine.Run(":7000")
	if err != nil {
		log.Fatal("gin engine coudn't start")
	}
}
