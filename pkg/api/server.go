package server

import (
	handlers "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/api/handler"
	middlewares "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/api/middleware"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartServer(
	user *handlers.UserHandler,
	admin *handlers.AdminHandler,
	cart *handlers.CartHandler,
	category *handlers.CategoryHandler,
	wishlist *handlers.WishlistHandler,
	product *handlers.ProductHandler,
	payment *handlers.PaymentHandler,
	order *handlers.OrderHandler,
	coupon *handlers.CouponHandler,
	adminMiddleware *middlewares.AdminMiddleware,
	userMiddleware *middlewares.UserMiddleware,
) {
	engine := gin.New()

	engine.Use(gin.Logger())

	engine.LoadHTMLFiles("/home/jasim/CityVibe-Ecommerce-CleanCode-Project/template/*")

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.AdminRoutes(engine, admin, product, category, order, coupon,adminMiddleware)
	routes.UserRoutes(engine, user, product, admin, wishlist, cart, order, payment, coupon,userMiddleware)

	engine.Run(":8080")
	
}
