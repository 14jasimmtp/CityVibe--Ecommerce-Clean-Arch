package usecase

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	interfaceRepo "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repository/interface"
	interfaceUsecase "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/usecase/interface"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type AdminUseCase struct {
	AdminRepo interfaceRepo.AdminRepo
	UserRepo interfaceRepo.UserRepo
	ProductRepo interfaceRepo.ProductRepo
	OrderRepo interfaceRepo.OrderRepo
}

func NewAdminUsecase(repo interfaceRepo.AdminRepo) interfaceUsecase.AdminUsecase {
	return &AdminUseCase{AdminRepo: repo}
}

func (clean *AdminUseCase) AdminLogin(admin models.AdminLogin) (models.Admin, error) {
	AdminDetails, err := clean.AdminRepo.AdminLogin(admin)
	fmt.Println(err)
	if err != nil {
		fmt.Println("Admin doesn't exist")
		return models.Admin{}, errors.New("admin not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(AdminDetails.Password), []byte(admin.Password)) != nil {
		fmt.Println("wrong password")
		return models.Admin{}, errors.New("wrong password")
	}

	tokenString, err := utils.AdminTokenGenerate(AdminDetails, "admin")
	if err != nil {
		fmt.Println("error generating token")
		return models.Admin{}, errors.New("error generating token")
	}

	return models.Admin{
		Firstname:   AdminDetails.Firstname,
		TokenString: tokenString,
	}, nil

}

func (clean *AdminUseCase) GetAllUsers() ([]models.UserDetailsResponse, error) {
	Users, err := clean.AdminRepo.GetAllUsers()
	if err != nil {
		return []models.UserDetailsResponse{}, err
	}
	return Users, nil
}

func (clean *AdminUseCase) BlockUser(idStr string) error {
	id, _ := strconv.Atoi(idStr)
	user, err := clean.UserRepo.GetUserById(id)
	if err != nil {
		return err
	}
	if user.Blocked {
		return errors.New("already blocked")
	} else {
		user.Blocked = true
	}
	err = clean.AdminRepo.BlockUserByID(user)
	if err != nil {
		return err
	}
	return nil

}

func (clean *AdminUseCase) UnBlockUser(idStr string) error {
	id, _ := strconv.Atoi(idStr)
	user, err := clean.UserRepo.GetUserById(id)
	if err != nil {
		return err
	}
	if !user.Blocked {
		return errors.New("already unblocked")
	} else {
		user.Blocked = false
	}
	err = clean.AdminRepo.UnBlockUserByID(user)
	if err != nil {
		return err
	}
	return nil

}

func (clean *AdminUseCase) GetAllOrderDetailsForAdmin() ([]models.ViewAdminOrderDetails, error) {
	orderDetail, err := clean.AdminRepo.GetAllOrderDetailsBrief()
	if err != nil {
		return []models.ViewAdminOrderDetails{}, err
	}
	return orderDetail, nil
}

func (clean *AdminUseCase) GetOrderDetails(orderID string) ([]models.OrderProductDetails, error) {
	orderDetails, err := clean.AdminRepo.GetSingleOrderDetails(orderID)
	if err != nil {
		return []models.OrderProductDetails{}, err
	}

	return orderDetails, nil
}

func (clean *AdminUseCase) ExecuteGetOffers() (*[]models.Offer, error) {
	offers, err := clean.ProductRepo.GetAllOffers()
	if err != nil {
		return nil, err
	}
	avialableoffers := []models.Offer{}
	for _, offers := range offers {
		if offers.UsageLimit != offers.UsedCount {
			avialableoffers = append(avialableoffers, offers)
		}
	}
	return &avialableoffers, nil
}

func (clean *AdminUseCase) ExecuteAddProductOffer(productid, offer int) (*models.Product, error) {

	product, err := clean.ProductRepo.GetProductById(productid)
	if err != nil {
		return nil, err
	}
	if offer < 0 || offer > 100 {
		return nil, errors.New("invalid offer percentage")
	}

	amount := float64(offer) / 100.0 * float64(product.Price)
	product.OfferPrize = math.Round((float64(product.Price)-amount)*100) / 100
	err1 := clean.ProductRepo.UpdateProduct(product)
	if err1 != nil {
		return nil, err1
	}
	return product, nil
}

func (clean *AdminUseCase) ExecuteCategoryOffer(catid, offer int) ([]models.Product, error) {

	productlist, err := clean.ProductRepo.GetProductsByCategoryoffer(catid)
	if err != nil {
		return nil, err
	}
	if offer < 0 || offer > 100 {
		return nil, errors.New("invalid offer percentage")
	}
	for i := range productlist {
		product := &(productlist)[i]

		amount := float64(offer) / 100.0 * float64(product.Price)
		product.OfferPrize = math.Round((float64(product.Price)-amount)*100) / 100
		err1 := clean.ProductRepo.UpdateProduct(product)
		if err1 != nil {
			return nil, err1
		}
	}
	return productlist, nil

}

func (clean *AdminUseCase) DashBoard() (models.CompleteAdminDashboard, error) {
	userDetails, err := clean.AdminRepo.DashBoardUserDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	productDetails, err := clean.AdminRepo.DashBoardProductDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	orderDetails, err := clean.OrderRepo.DashBoardOrder()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	totalRevenue, err := clean.AdminRepo.TotalRevenue()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	amountDetails, err := clean.AdminRepo.AmountDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	return models.CompleteAdminDashboard{
		DashboardUser:    userDetails,
		DashboardProduct: productDetails,
		DashboardOrder:   orderDetails,
		DashboardRevenue: totalRevenue,
		DashboardAmount:  amountDetails,
	}, nil
}

func (clean *AdminUseCase) ExecuteAddOffer(offer *models.Offer) error {
	err := clean.AdminRepo.CreateOffer(offer)
	if err != nil {
		return errors.New("error creating offer")
	} else {
		return nil
	}
}

