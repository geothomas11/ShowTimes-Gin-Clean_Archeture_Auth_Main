package routes

import (
	"ShowTimes/pkg/api/handler"
	"ShowTimes/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.RouterGroup, userHandler *handler.UserHandler, otphandler *handler.OtpHandler, cartHandler *handler.CartHandler, orderHandler *handler.OrderHandler, ProductHandler *handler.ProductHandler, paymentHandler *handler.PaymentHandler, walletHandler *handler.WalletHandler, couponHandler *handler.CouponHandler) {
	engine.GET("/google_callback", userHandler.GoogleCallback)
	engine.GET("/google_login", userHandler.Authv2)

	engine.POST("/signup", userHandler.UserSignUp)
	engine.POST("/login", userHandler.LoginHandler)

	engine.GET("/products", ProductHandler.ListProductsUser)

	engine.POST("/otplogin", otphandler.SendOTP)
	engine.POST("/verifyotp", otphandler.VerifyOTP)

	engine.GET("/payment", paymentHandler.MakePaymentRazorpay)
	engine.GET("/verifypayment", paymentHandler.VerifyPayment)

	engine.Use(middleware.UserAuthMiddleware)
	{

		profile := engine.Group("/profile")
		{
			profile.POST("/addaddress", userHandler.AddAddress)
			profile.GET("/showuserdetails", userHandler.ShowUserDetails)
			profile.GET("/alladdress", userHandler.GetAllAddress)
			profile.PUT("/editprofile", userHandler.EditProfile)
			profile.PATCH("/changepassword", userHandler.ChangePassword)
		}
		cart := engine.Group("/cart")
		{
			cart.POST("/addtocart", cartHandler.AddToCart)
			cart.GET("/listcartitems", cartHandler.ListCartItems)
			cart.PATCH("/updateproductquantity", cartHandler.UpdateProductQuantityCart)
			cart.PUT("/removefromcart", cartHandler.RemoveFromCart)

		}

		Checkout := engine.Group("/orders")
		{
			Checkout.GET("/checkout", orderHandler.Checkout)
			Checkout.POST("/Orderitemsfromcart", orderHandler.OrderItems)
			Checkout.GET("/orderDetails", orderHandler.GetOrderDetails)
			Checkout.DELETE("/cancelOrder", orderHandler.CancelOrder)
			// Checkout.POST("/placeorderinCOD", orderHandler.PlaceOrderCOD)
			Checkout.PATCH("/returnorder", orderHandler.ReturnOrder)
			Checkout.GET("/print", orderHandler.PrintInvoice)

		}
		wallet := engine.Group("/wallet", walletHandler.GetWallet)
		{
			wallet.GET("/getwallet", walletHandler.GetWallet)
		}
		coupon := engine.Group("/coupon")
		{
			coupon.GET("/addcoupon", couponHandler.GetCouponUser)
		}

	}
}
