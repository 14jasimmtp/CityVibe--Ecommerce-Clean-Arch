package di

import (
	server "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/api"
	handlers "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/api/handler"
	middlewares "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/api/middleware"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/config"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/db"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repository"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/usecase"
)

func InitialiseAPI(cfg config.Config) {
	DB := db.DBInitialise(cfg)
	
	userMiddleware := middlewares.NewUserMiddleware()
	adminMiddleware := middlewares.NewAdminMiddleware()
 
	userRepository := repository.NewUserRepo(DB)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userHandler := handlers.NewUserHandler(userUsecase)

	productRepository := repository.NewProductRepo(DB)
	productUsecase := usecase.NewProductUsecase(productRepository)
	productHandler := handlers.NewProductHandler(productUsecase)
	
	cartRepository := repository.NewCartRepo(DB)
	cartUsecase := usecase.NewCartUsecase(cartRepository)
	cartHandler := handlers.NewCartHandler(cartUsecase)

	categoryRepository := repository.NewCategoryRepo(DB)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryUsecase)

	couponRepository := repository.NewCouponRepo(DB)
	couponUsecase := usecase.NewCouponUsecase(couponRepository)
	couponHandler := handlers.NewCouponHandler(couponUsecase)

	paymentRepository := repository.NewPaymentRepo(DB)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepository)
	paymentHandler := handlers.NewPaymentHandler(paymentUsecase)

	orderRepository := repository.NewOrderRepo(DB)
	orderUsecase := usecase.NewOrderUsecase(orderRepository,couponRepository,userRepository,paymentRepository,cartRepository)
	orderHandler := handlers.NewOrderHandler(orderUsecase)

	adminRepository := repository.NewAdminRepo(DB)
	adminUsecase := usecase.NewAdminUsecase(adminRepository,userRepository,productRepository,orderRepository)
	adminHandler := handlers.NewAdminHandler(adminUsecase)

	wishlistRepository := repository.NewWishlistRepo(DB)
	wishlistUsecase := usecase.NewWishlistUsecase(wishlistRepository)
	wishlistHandler := handlers.NewWishlistHandler(wishlistUsecase)

	server.StartServer(
		userHandler,
		adminHandler,
		cartHandler,
		categoryHandler,
		wishlistHandler,
		productHandler,
		paymentHandler,
		orderHandler,
		couponHandler,
		adminMiddleware,
		userMiddleware,
	)
}