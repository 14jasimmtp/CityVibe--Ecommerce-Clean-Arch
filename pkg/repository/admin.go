package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/db"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	interfaceRepo "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repository/interface"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/utils"
	"gorm.io/gorm"
)

type AdminRepo struct {
	DB *gorm.DB
}

func NewAdminRepo(db *gorm.DB) interfaceRepo.AdminRepo{
	return &AdminRepo{DB: db}
}

func (clean *AdminRepo) AdminLogin(adminDetails models.AdminLogin) (models.Admin, error) {
	var details models.Admin
	if err := db.DB.Raw("SELECT * FROM admins WHERE email=?", adminDetails.Email).Scan(&details).Error; err != nil {
		return models.Admin{}, err
	}
	return details, nil
}

func (clean *AdminRepo) GetAllUsers() ([]models.UserDetailsResponse, error) {
	var users []models.UserDetailsResponse
	result := db.DB.Raw("SELECT id,email,firstname,lastname,phone,blocked,wallet FROM users").Scan(&users)
	if result.Error != nil {
		fmt.Println("data fetching error")
		return []models.UserDetailsResponse{}, result.Error
	}

	return users, nil
}

func (clean *AdminRepo) BlockUserByID(user models.UserDetailsResponse) error {
	result := db.DB.Exec("UPDATE users SET blocked = ? WHERE id = ?", user.Blocked, user.ID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (clean *AdminRepo) UnBlockUserByID(user models.UserDetailsResponse) error {
	result := db.DB.Exec("UPDATE users SET blocked = ? WHERE id = ?", user.Blocked, user.ID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (clean *AdminRepo) GetAllOrderDetailsBrief() ([]models.ViewAdminOrderDetails, error) {

	var orderDatails []models.AdminOrderDetails
	query := db.DB.Raw("SELECT orders.user_id,orders.id, total_price as final_price, payment_methods.payment_mode AS payment_method, payment_status FROM orders INNER JOIN payment_methods ON orders.payment_method_id=payment_methods.id  ORDER BY orders.id DESC").Scan(&orderDatails)
	if query.Error != nil {
		return []models.ViewAdminOrderDetails{}, errors.New(`something went wrong`)
	}
	var fullOrderDetails []models.ViewAdminOrderDetails
	for _, ok := range orderDatails {
		var OrderProductDetails []models.OrderProductDetails
		db.DB.Raw("SELECT order_items.product_id,products.name AS product_name,order_items.order_status,order_items.quantity,order_items.total_price FROM order_items INNER JOIN products ON order_items.product_id = products.id WHERE order_items.order_id = ? ORDER BY order_id DESC", ok.Id).Scan(&OrderProductDetails)
		fullOrderDetails = append(fullOrderDetails, models.ViewAdminOrderDetails{OrderDetails: ok, OrderProductDetails: OrderProductDetails})
	}
	return fullOrderDetails, nil

}

func (clean *AdminRepo) GetSingleOrderDetails(orderID string) ([]models.OrderProductDetails, error) {
	var Order []models.OrderProductDetails
	query := db.DB.Raw(`SELECT product_id,products.name AS product_name,order_status,quantity,Total_price FROM order_items INNER JOIN products ON product_id=products.id WHERE order_id = ?`, orderID).Scan(&Order)
	if query.Error != nil {
		return []models.OrderProductDetails{}, query.Error
	}
	return Order, nil
}

func (clean *AdminRepo) DashBoardUserDetails() (models.DashBoardUser, error) {
	var userDetails models.DashBoardUser
	err := db.DB.Raw("SELECT COUNT(*) FROM users").Scan(&userDetails.TotalUsers).Error
	if err != nil {
		return models.DashBoardUser{}, nil
	}
	err = db.DB.Raw("SELECT id FROM users WHERE blocked=true").Scan(&userDetails.BlockedUser).Error
	if err != nil {
		return models.DashBoardUser{}, nil
	}
	return userDetails, nil
}

func (clean *AdminRepo) DashBoardProductDetails() (models.DashBoardProduct, error) {
	var productDetails models.DashBoardProduct
	err := db.DB.Raw("SELECT COUNT(*) FROM products").Scan(&productDetails.TotalProducts).Error
	if err != nil {
		return models.DashBoardProduct{}, nil
	}
	err = db.DB.Raw("SELECT id FROM products WHERE stock=0").Scan(&productDetails.OutofStockProductID).Error
	if err != nil {
		return models.DashBoardProduct{}, nil
	}
	err = db.DB.Raw("SELECT id FROM products WHERE stock<=5").Scan(&productDetails.LowStockProductsID).Error
	if err != nil {
		return models.DashBoardProduct{}, nil
	}
	return productDetails, nil
}

func (clean *AdminRepo) TotalRevenue() (models.DashboardRevenue, error) {
	var revenueDetails models.DashboardRevenue
	startTime := time.Now().AddDate(0, 0, -1)
	endTime := time.Now()
	err := db.DB.Raw("SELECT COALESCE(SUM(total_price),0) FROM orders WHERE payment_status = 'paid' AND created_at >=? AND created_at <=?", startTime, endTime).Scan(&revenueDetails.TodayRevenue).Error
	if err != nil {
		return models.DashboardRevenue{}, nil
	}
	startTime, endTime = utils.CalcualtePeriodDate("monthly")
	err = db.DB.Raw("SELECT COALESCE (SUM(total_price),0) FROM orders WHERE payment_status = 'paid' AND created_at >=? AND created_at <=?", startTime, endTime).Scan(&revenueDetails.MonthRevenue).Error
	if err != nil {
		return models.DashboardRevenue{}, nil
	}
	startTime, endTime = utils.CalcualtePeriodDate("yearly")
	err = db.DB.Raw("SELECT COALESCE (SUM(total_price),0) FROM orders WHERE payment_status = 'paid' AND created_at >=? AND created_at <=?", startTime, endTime).Scan(&revenueDetails.YearRevenue).Error
	if err != nil {
		return models.DashboardRevenue{}, nil
	}
	return revenueDetails, nil
}

func (clean *AdminRepo) AmountDetails() (models.DashboardAmount, error) {
	var amountDetails models.DashboardAmount
	err := db.DB.Raw("SELECT COALESCE (SUM(total_price),0) FROM orders WHERE payment_status = 'paid' ").Scan(&amountDetails.CreditedAmount).Error
	if err != nil {
		return models.DashboardAmount{}, nil
	}
	err = db.DB.Raw("SELECT COALESCE(SUM(total_price),0) FROM orders WHERE payment_status = 'not paid' ").Scan(&amountDetails.PendingAmount).Error
	if err != nil {
		return models.DashboardAmount{}, nil
	}
	return amountDetails, nil
}

func (clean *AdminRepo) CreateOffer(offer *models.Offer) error {
	if err := db.DB.Create(offer).Error; err != nil {
		return err
	}
	return nil
}
