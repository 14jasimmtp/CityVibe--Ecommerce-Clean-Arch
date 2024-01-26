package interfaceUsecase

import "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"

type AdminUsecase interface {
	AdminLogin(admin models.AdminLogin) (models.Admin, error)
	GetAllUsers() ([]models.UserDetailsResponse, error)
	BlockUser(idStr string) error
	UnBlockUser(idStr string) error
	GetAllOrderDetailsForAdmin() ([]models.ViewAdminOrderDetails, error)
	GetOrderDetails(orderID string) ([]models.OrderProductDetails, error)
	ExecuteGetOffers() (*[]models.Offer, error)
	ExecuteAddProductOffer(productid, offer int) (*models.Product, error)
	ExecuteCategoryOffer(catid, offer int) ([]models.Product, error)
	DashBoard() (models.CompleteAdminDashboard, error)
	ExecuteAddOffer(offer *models.Offer) error
}
