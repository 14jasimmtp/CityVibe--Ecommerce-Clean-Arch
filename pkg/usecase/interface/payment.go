package interfaceUsecase

import "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"

type PaymentUsecase interface {
	MakePaymentRazorPay(orderID int) (models.Payment, string, error)
	PaymentMethodID(orderID int) (int, error)
	PaymentAlreadyPaid(orderID int) (bool, error)
	VerifyPayment(details models.PaymentVerify, orderID int) (models.OrderDetails, error)
}
