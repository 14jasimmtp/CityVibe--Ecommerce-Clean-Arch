package interfaceRepo

import "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"

type AdminRepo interface {
	AdminLogin(adminDetails models.AdminLogin) (models.Admin, error)
	GetAllUsers() ([]models.UserDetailsResponse, error)
	BlockUserByID(user models.UserDetailsResponse) error
	UnBlockUserByID(user models.UserDetailsResponse) error
	GetAllOrderDetailsBrief() ([]models.ViewAdminOrderDetails, error)
	GetSingleOrderDetails(orderID string) ([]models.OrderProductDetails, error)
	DashBoardUserDetails() (models.DashBoardUser, error)
	DashBoardProductDetails() (models.DashBoardProduct, error)
	TotalRevenue() (models.DashboardRevenue, error)
	AmountDetails() (models.DashboardAmount, error)
	CreateOffer(offer *models.Offer) error
}
