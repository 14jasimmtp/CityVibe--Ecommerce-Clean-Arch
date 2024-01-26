package routes

import (
	handlers "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/api/handler"
	middlewares "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/api/middleware"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(
	r *gin.Engine,
	admin *handlers.AdminHandler,
	product *handlers.ProductHandler,
	category *handlers.CategoryHandler,
	order *handlers.OrderHandler,
	coupon *handlers.CouponHandler,
	adminMiddleware *middlewares.AdminMiddleware,
) {

	//USER
	r.POST("/admin/login", admin.AdminLogin)
	r.GET("/admin/users", adminMiddleware.AdminAuthMiddleware, admin.GetAllUsers)
	r.POST("/admin/users/block", adminMiddleware.AdminAuthMiddleware, admin.BlockUser)
	r.POST("/admin/users/unblock", adminMiddleware.AdminAuthMiddleware, admin.UnBlockUser)

	//product
	r.GET("/admin/products", adminMiddleware.AdminAuthMiddleware, product.AllProducts)
	r.POST("/admin/products", adminMiddleware.AdminAuthMiddleware, product.AddProduct)
	r.PUT("/admin/products", adminMiddleware.AdminAuthMiddleware, product.EditProductDetails)
	r.DELETE("admin/products/remove/:id", adminMiddleware.AdminAuthMiddleware, product.DeleteProduct)

	//category
	r.GET("admin/category", adminMiddleware.AdminAuthMiddleware, category.GetCategory)
	r.POST("admin/category", adminMiddleware.AdminAuthMiddleware, category.AddCategory)
	r.PUT("admin/category", adminMiddleware.AdminAuthMiddleware, category.UpdateCategory)
	r.DELETE("admin/category", adminMiddleware.AdminAuthMiddleware, category.DeleteCategory)

	//order
	r.GET("admin/orders", adminMiddleware.AdminAuthMiddleware, admin.OrderDetailsForAdmin)
	r.POST("admin/orders/ship", adminMiddleware.AdminAuthMiddleware, order.ShipOrderByAdmin)
	r.POST("admin/orders/cancel", adminMiddleware.AdminAuthMiddleware, order.CancelOrderByAdmin)
	r.GET("admin/orders/details", adminMiddleware.AdminAuthMiddleware, admin.OrderDetailsforAdminWithID)
	r.POST("admin/orders/deliver", adminMiddleware.AdminAuthMiddleware, order.DeliverOrderByAdmin)

	//coupons
	r.POST("admin/coupon", adminMiddleware.AdminAuthMiddleware, coupon.MakeCoupon)
	r.PUT("admin/coupon/disable", adminMiddleware.AdminAuthMiddleware, coupon.DisableCoupon)
	r.PUT("admin/coupon/enable", adminMiddleware.AdminAuthMiddleware, coupon.EnableCoupon)
	r.GET("admin/coupon", adminMiddleware.AdminAuthMiddleware, coupon.ViewCouponsAdmin)
	r.PUT("admin/coupon/update", adminMiddleware.AdminAuthMiddleware, coupon.UpdateCoupon)

	//salesreport
	r.GET("admin/salesreportbyperiod", adminMiddleware.AdminAuthMiddleware, order.SalesReportByPeriod)
	r.GET("admin/salesreportbydate", adminMiddleware.AdminAuthMiddleware, order.SalesReportByDate)
	r.GET("admin/salesreportbypayment", adminMiddleware.AdminAuthMiddleware, order.SalesReportByPayment)
	r.POST("admin/salesreport/excel", adminMiddleware.AdminAuthMiddleware, order.SalesReportXL)

	//dashboard
	r.GET("/admin/dashboard", adminMiddleware.AdminAuthMiddleware, admin.DashBoard)

	//offer
	r.POST("admin/product/offer", adminMiddleware.AdminAuthMiddleware, admin.AddProductOffer)
	r.POST("admin/category/offer", adminMiddleware.AdminAuthMiddleware, admin.AddCategoryOffer)
}
