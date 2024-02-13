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
	DB:=db.DBInitialise(cfg)

	userMiddleware:=middlewares.NewUserMiddleware()
	adminMiddleware:=middlewares.NewAdminMiddleware()

	adminRepository:=repository.NewAdminRepo(DB)
	adminUsecase:=usecase.NewAdminUsecase(adminRepository)
	adminHandler:=handlers.NewAdminHandler(adminUsecase)

	userRepository:=repository.NewUserRepo(DB)
	userUsecase:=usecase.NewUserUsecase(userRepository)
	userHandler:=handlers.NewUserHandler(userUsecase)

	cartRepository:=repository.NewCartRepo(DB)
	cartUsecase:=usecase.NewCartUsecase(cartRepository)
	cartHandler:=handlers.NewCartHandler(cartUsecase)

	categoryRepository:=repository.NewCategoryRepo(DB)
	categoryUsecase:=usecase.NewCategoryUsecase(categoryRepository)
	categoryHandler:=handlers.NewCategoryHandler(categoryUsecase)

	couponRepository:=repository.NewCouponRepo(DB)
	couponUsecase:=usecase.NewCouponUsecase(couponRepository)
	couponHandler:=handlers.NewCouponHandler(couponUsecase)

	orderRepository:=repository.NewOrderRepo(DB)
	orderUsecase:=usecase.NewOrderUsecase(orderRepository)
	orderHandler:=handlers.NewOrderHandler(orderUsecase)

	paymentRepository:=repository.NewPaymentRepo(DB)
	paymentUsecase:=usecase.NewPaymentUsecase(paymentRepository)
	paymentHandler:=handlers.NewPaymentHandler(paymentUsecase)

	productRepository:=repository.NewProductRepo(DB)
	productUsecase:=usecase.NewProductUsecase(productRepository)
	productHandler:=handlers.NewProductHandler(productUsecase)

	wishlistRepository:=repository.NewWishlistRepo(DB)
	wishlistUsecase:=usecase.NewWishlistUsecase(wishlistRepository)
	wishlistHandler:=handlers.NewWishlistHandler(wishlistUsecase)

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