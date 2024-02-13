package usecase

import (
	"errors"
	"fmt"
	"log"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/config"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	interfaceRepo "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repository/interface"
	interfaceUsecase "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/usecase/interface"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/utils"
	"github.com/razorpay/razorpay-go"
)

type PaymentUseCase struct {
	PaymentRepo interfaceRepo.PaymentRepo
}

func NewPaymentUsecase(repo interfaceRepo.PaymentRepo) interfaceUsecase.PaymentUsecase {
	return &PaymentUseCase{PaymentRepo: repo}
}

func (clean *PaymentUseCase) MakePaymentRazorPay(orderID int) (models.Payment, string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	PaymentDetails, err := clean.PaymentRepo.GetPaymentDetails(orderID)
	if err != nil {
		return models.Payment{}, "", err
	}

	client := razorpay.NewClient(cfg.KEY_ID_FOR_PAY, cfg.SECRET_KEY_FOR_PAY)

	data := map[string]interface{}{
		"amount":   int(PaymentDetails.Final_price * 100),
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		fmt.Println("hello")
		fmt.Println(err)
		return models.Payment{}, "", err
	}

	razorPayOrderID := body["id"].(string)

	err = clean.PaymentRepo.AddRazorPayDetails(orderID, razorPayOrderID)
	if err != nil {
		fmt.Println("hig")
		return models.Payment{}, "", err
	}

	return PaymentDetails, razorPayOrderID, nil

}

func (clean *PaymentUseCase) PaymentMethodID(orderID int) (int, error) {
	PaymethodID, err := clean.PaymentRepo.PayMethod(orderID)
	if err != nil {
		return 0, err
	}
	return PaymethodID, nil
}

func (clean *PaymentUseCase) PaymentAlreadyPaid(orderID int) (bool, error) {
	AlreadyPayed, err := clean.PaymentRepo.PaymentAlreadyPaid(orderID)
	if err != nil {
		return false, err
	}
	return AlreadyPayed, nil
}

func (clean *PaymentUseCase) VerifyPayment(details models.PaymentVerify, order_id int) (models.OrderDetails, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	paid, err := clean.PaymentRepo.CheckVerifiedPayment(order_id)
	if err != nil {
		return models.OrderDetails{}, err
	}
	if paid {
		return models.OrderDetails{}, errors.New(`already payment verified`)
	}

	result := utils.VerifyPayment(details.OrderID, details.PaymentID, details.Signature, cfg.SECRET_KEY_FOR_PAY)
	if !result {
		return models.OrderDetails{}, errors.New("payment is unsuccessful")
	}

	orders, err := clean.PaymentRepo.UpdateShipmentAndPaymentByOrderID("processing", "paid", order_id)
	if err != nil {
		return models.OrderDetails{}, err
	}
	return orders, nil
}
