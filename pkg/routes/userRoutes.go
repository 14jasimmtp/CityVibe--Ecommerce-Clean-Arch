package routes

import (
	handlers "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/api/handler"
	middlewares "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/api/middleware"
	"github.com/gin-gonic/gin"
)



func UserRoutes(
	r *gin.Engine,
	user *handlers.UserHandler,
	product *handlers.ProductHandler,
	admin *handlers.AdminHandler,
	wishlist *handlers.WishlistHandler,
	cart *handlers.CartHandler,
	order *handlers.OrderHandler,
	razorpay *handlers.PaymentHandler,
	coupon *handlers.CouponHandler,
	userMiddleware *middlewares.UserMiddleware,
 ) {
	//login
	r.POST("/signup", user.UserSignup)
	r.POST("/login", user.UserLogin)

	r.POST("/verify", user.VerifyLoginOtp)
	r.POST("/logout", user.UserLogout)

	//products
	r.GET("/products", product.GetAllProducts)
	r.GET("/products/:id", product.ShowSingleProduct)
	r.GET("/products/search", product.SearchProducts) //search

	//filtering
	r.GET("/products/filter", product.FilterProducts)

	//wishlist
	r.POST("/products/wishlist", userMiddleware.UserAuthMiddleware, wishlist.AddProductToWishlist)
	r.GET("/products/wishlist", userMiddleware.UserAuthMiddleware, wishlist.ViewUserWishlist)
	r.DELETE("/products/wishlist", userMiddleware.UserAuthMiddleware, wishlist.RemoveProductFromWishlist)

	//profile
	r.GET("/profile", userMiddleware.UserAuthMiddleware, user.UserProfile)
	r.PUT("/profile", userMiddleware.UserAuthMiddleware, user.UpdateUserProfile)

	//password change
	r.POST("/password/forgot", user.ForgotPassword)
	r.POST("password/forgot/change", user.ResetForgottenPassword)

	//Address
	r.GET("/address", userMiddleware.UserAuthMiddleware, user.ViewUserAddress)
	r.POST("/address", userMiddleware.UserAuthMiddleware, user.AddNewAddressDetails)
	r.PUT("/address", userMiddleware.UserAuthMiddleware, user.EditUserAddress)
	r.DELETE("/address", userMiddleware.UserAuthMiddleware, user.RemoveUserAddress)

	//cart
	r.GET("/cart", userMiddleware.UserAuthMiddleware, cart.ViewCart)
	r.POST("/cart", userMiddleware.UserAuthMiddleware, cart.AddToCart)
	r.DELETE("/cart", userMiddleware.UserAuthMiddleware, cart.RemoveProductsFromCart)
	r.PUT("/cart/add-quantity", userMiddleware.UserAuthMiddleware, cart.IncreaseQuantityUpdate)
	r.PUT("/cart/reduce-quantity", userMiddleware.UserAuthMiddleware, cart.DecreaseQuantityUpdate)

	//orders
	r.GET("/orders", userMiddleware.UserAuthMiddleware, order.ViewOrders)
	r.POST("/orders", userMiddleware.UserAuthMiddleware, order.OrderFromCart)
	r.GET("/checkout", userMiddleware.UserAuthMiddleware, order.ViewCheckOut)
	r.PUT("/orders/return", userMiddleware.UserAuthMiddleware, order.ReturnOrder)
	r.PUT("/orders/cancel", userMiddleware.UserAuthMiddleware, order.CancelOrder)

	//wishlist
	r.GET("/wishlist", userMiddleware.UserAuthMiddleware, wishlist.ViewUserWishlist)
	r.POST("/wishlist", userMiddleware.UserAuthMiddleware, wishlist.AddProductToWishlist)
	r.DELETE("/wishlist", userMiddleware.UserAuthMiddleware, wishlist.RemoveProductFromWishlist)

	//payment
	r.GET("/payment/razorpay", razorpay.ExecuteRazorPayPayment)
	r.POST("/payment/verify", razorpay.VerifyPayment)

	//coupons
	r.GET("/coupons", userMiddleware.UserAuthMiddleware, coupon.ViewCouponsUser)
	r.POST("/applycoupon",userMiddleware.UserAuthMiddleware,order.ApplyCoupon)
	r.POST("/removecoupon",userMiddleware.UserAuthMiddleware,coupon.RemoveCoupon)

	//Invoice
	r.GET("/Invoice", userMiddleware.UserAuthMiddleware, order.PrintInvoice)

}
