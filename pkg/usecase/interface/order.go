package interfaceUsecase

import (
	"time"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	"github.com/xuri/excelize/v2"
	"github.com/jung-kurt/gofpdf"
)

type OrderUseCase interface {
	CheckOut(Token string) (interface{}, error)
	ExecutePurchase(Token string, OrderInput models.CheckOut) (models.OrderSuccessResponse, error)
	ExecutePurchaseWallet(Token string, OrderInput models.CheckOut) (models.OrderSuccessResponse, error)
	ViewUserOrders(Token string) ([]models.ViewOrderDetails, error)
	CancelOrder(Token, orderID, pid string) error
	CancelOrderByAdmin(userID, orderID, pid int) error
	ShipOrders(userID, orderId, pid int) error
	DeliverOrder(userID, orderId, pid int) error
	ReturnOrder(Token, orderID, pid string) error
	ExecuteSalesReportByPeriod(period string) (*gofpdf.Fpdf, error)
	ExecuteSalesReportByDate(startdate, enddate time.Time) (*gofpdf.Fpdf, error)
	ExecuteSalesReportByPaymentMethod(startdate, enddate time.Time, paymentmethod string) (*gofpdf.Fpdf, error)
	PrintInvoice(orderID int, Token string) (*gofpdf.Fpdf, error)
	ApplyCoupon(coupon, Token string) error
	SalesReportXL(start, end time.Time) (*excelize.File, error)
}
