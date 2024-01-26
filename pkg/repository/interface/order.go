package interfaceRepo

import (
	"time"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
)

type OrderRepo interface {
	OrderFromCart(addressid uint, paymentid, userid uint, price float64) (int, error)
	AddOrderProducts(userID uint, orderid int, cart []models.Cart) error
	CheckPaymentMethodExist(paymentid uint) bool
	GetOrder(orderID int) (domain.Order, error)
	GetOrderDetails(userID uint) ([]models.ViewOrderDetails, error)
	CheckOrder(orderid string, userID uint) error
	GetProductDetailsFromOrders(orderid int) ([]models.Product, error)
	GetOrderStatus(orderId, pid string) (string, error)
	CancelOrder(orderid, pid string, userID uint) error
	UpdateStock(pid int, quantity int) error
	UpdateSingleStock(pid string) error
	UpdateCartAndStockAfterOrder(userID uint, productID int, quantity float64) error
	CheckSingleOrder(pid, orderId string, userId uint) error
	CancelSingleOrder(pid, orderId string, userId uint) error
	CancelOrderByAdmin(orderID string) error
	ShipOrder(userID, orderId int) error
	DeliverOrder(userID int, orderId string) error
	UpdateFinalPrice(userID uint, oid string) error
	ReturnAmountToWallet(userID uint, orderID, pid string) error
	CancelOrderDetails(userID uint, orderID, pid string) (models.CancelDetails, error)
	UpdateOrderFinalPrice(orderID int, amount float64) error
	UpdateCartAmount(userID, discount uint) (float64, error)
	ReturnOrder(userID uint, orderID, pid string) error
	GetOrderInvoice(orderID, UserID int) (domain.Order, error)
	GetByDate(startdate, enddate time.Time) (*models.SalesReport, error)
	GetByPaymentMethod(startdate, enddate time.Time, paymentmethod string) (*models.SalesReport, error)
	GetAddressFromOrders(address_id, userID int) (models.Address, error)
	DashBoardOrder() (models.DashboardOrder, error)
	XLBYDATE(start, end time.Time) ([]models.SalesReportXL, error)
}
