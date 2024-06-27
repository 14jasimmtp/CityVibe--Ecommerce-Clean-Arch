package usecase

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	interfaceRepo "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repository/interface"
	interfaceUsecase "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/usecase/interface"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/utils"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type OrderUseCase struct {
	OrderRepo   interfaceRepo.OrderRepo
	CouponRepo  interfaceRepo.CouponRepo
	UserRepo    interfaceRepo.UserRepo
	PaymentRepo interfaceRepo.PaymentRepo
	CartRepo    interfaceRepo.CartRepo
}

func NewOrderUsecase(Orderrepo interfaceRepo.OrderRepo,Couponrepo interfaceRepo.CouponRepo,Userrepo interfaceRepo.UserRepo,Paymentrepo interfaceRepo.PaymentRepo, cartrepo interfaceRepo.CartRepo) interfaceUsecase.OrderUseCase {
	return &OrderUseCase{OrderRepo: Orderrepo,CouponRepo: Couponrepo,UserRepo: Userrepo,PaymentRepo: Paymentrepo,CartRepo: cartrepo}
}

func (clean *OrderUseCase) CheckOut(Token string) (interface{}, error) {
	userId, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return models.CheckOutInfo{}, err
	}

	AllUserAddress, err := clean.UserRepo.ViewAddress(userId)
	if err != nil {
		return models.CheckOutInfo{}, err
	}

	AllCartProducts, err := clean.CartRepo.DisplayCart(userId)
	if err != nil {
		return models.CheckOutInfo{}, err
	}

	TotalAmount, err := clean.CartRepo.CartTotalAmount(userId)
	if err != nil {
		return models.CheckOutInfo{}, err
	}

	if AllCartProducts[0].FinalPrice == 0 {
		return models.CheckOutInfo{
			Address:     AllUserAddress,
			Cart:        AllCartProducts,
			TotalAmount: TotalAmount,
		}, nil
	} else {

		finalPrice, err := clean.CartRepo.CartFinalPrice(userId)
		if err != nil {
			return models.CheckOutInfo{}, err
		}
		for i := 0; i < len(AllCartProducts); i++ {
			AllCartProducts[i].Price = AllCartProducts[i].FinalPrice
		}
		return models.CheckOutInfoDiscount{
			Address:        AllUserAddress,
			Cart:           AllCartProducts,
			TotalAmount:    TotalAmount,
			DiscountAmount: finalPrice,
		}, nil
	}
}

func (clean *OrderUseCase) ExecutePurchase(Token string, OrderInput models.CheckOut) (models.OrderSuccessResponse, error) {
	var TotalAmount float64
	var method string
	userId, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	addressExist := clean.UserRepo.CheckAddressExist(userId, OrderInput.AddressID)
	if !addressExist {
		return models.OrderSuccessResponse{}, errors.New(`address doesn't exist`)
	}

	paymentExist := clean.OrderRepo.CheckPaymentMethodExist(OrderInput.PaymentID)
	if !paymentExist {
		return models.OrderSuccessResponse{}, errors.New(`payment method doesn't exist`)
	}
	if OrderInput.PaymentID == 1 {
		method = "COD"
	} else {
		method = "Razorpay"
	}

	cartExist := clean.CartRepo.CheckCartExist(userId)
	if !cartExist {
		return models.OrderSuccessResponse{}, errors.New(`cart is empty`)
	}

	cartItems, err := clean.CartRepo.DisplayCart(userId)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	if cartItems[0].FinalPrice != 0 {

		for i := 0; i < len(cartItems); i++ {
			cartItems[i].Price = cartItems[i].FinalPrice
		}

		TotalAmount, err = clean.CartRepo.CartFinalPrice(userId)
		if err != nil {
			return models.OrderSuccessResponse{}, errors.New(`error while calculating total amount`)
		}
	} else {
		TotalAmount, err = clean.CartRepo.CartTotalAmount(userId)
		if err != nil {
			return models.OrderSuccessResponse{}, errors.New(`error while calculating total amount`)
		}
	}

	OrderID, err := clean.OrderRepo.OrderFromCart(OrderInput.AddressID, OrderInput.PaymentID, userId, TotalAmount)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	if err := clean.OrderRepo.AddOrderProducts(userId, OrderID, cartItems); err != nil {
		return models.OrderSuccessResponse{}, err
	}

	var orderItemDetails domain.OrderItem
	for _, c := range cartItems {
		orderItemDetails.ProductID = c.ProductID
		orderItemDetails.Quantity = c.Quantity
		err := clean.OrderRepo.UpdateCartAndStockAfterOrder(userId, int(orderItemDetails.ProductID), orderItemDetails.Quantity)
		if err != nil {
			return models.OrderSuccessResponse{}, err
		}
	}
	return models.OrderSuccessResponse{
		OrderID:       OrderID,
		PaymentMethod: method,
		TotalAmount:   TotalAmount,
		PaymentStatus: "not paid",
	}, nil
}

func (clean *OrderUseCase) ExecutePurchaseWallet(Token string, OrderInput models.CheckOut) (models.OrderSuccessResponse, error) {
	var TotalAmount float64
	userId, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	user, err := clean.UserRepo.GetUserById(int(userId))
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	addressExist := clean.UserRepo.CheckAddressExist(userId, OrderInput.AddressID)
	if !addressExist {
		return models.OrderSuccessResponse{}, errors.New(`address doesn't exist`)
	}

	paymentExist := clean.OrderRepo.CheckPaymentMethodExist(OrderInput.PaymentID)
	if !paymentExist {
		return models.OrderSuccessResponse{}, errors.New(`payment method doesn't exist`)
	}

	cartExist := clean.CartRepo.CheckCartExist(userId)
	if !cartExist {
		return models.OrderSuccessResponse{}, errors.New(`cart is empty`)
	}

	cartItems, err := clean.CartRepo.DisplayCart(userId)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	if cartItems[1].FinalPrice != 0 {

		for i := 0; i < len(cartItems); i++ {
			cartItems[i].Price = cartItems[i].FinalPrice
		}

		TotalAmount, err = clean.CartRepo.CartFinalPrice(userId)
		if err != nil {
			return models.OrderSuccessResponse{}, errors.New(`error while calculating total amount`)
		}
	} else {
		TotalAmount, err = clean.CartRepo.CartTotalAmount(userId)
		if err != nil {
			return models.OrderSuccessResponse{}, errors.New(`error while calculating total amount`)
		}
	}

	if user.Wallet < TotalAmount {
		return models.OrderSuccessResponse{}, errors.New(`insufficient Balance in Wallet.Add amount to wallet to purchase`)
	}

	OrderID, err := clean.OrderRepo.OrderFromCart(OrderInput.AddressID, OrderInput.PaymentID, userId, TotalAmount)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	if err := clean.OrderRepo.AddOrderProducts(userId, OrderID, cartItems); err != nil {
		return models.OrderSuccessResponse{}, err
	}
	_, err = clean.PaymentRepo.UpdateShipmentAndPaymentByOrderID("processing", "paid", OrderID)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	user.Wallet -= TotalAmount

	err = clean.UserRepo.UpdateWallet(user.Wallet, userId)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	var orderItemDetails domain.OrderItem
	for _, c := range cartItems {
		orderItemDetails.ProductID = c.ProductID
		orderItemDetails.Quantity = c.Quantity
		err := clean.OrderRepo.UpdateCartAndStockAfterOrder(userId, int(orderItemDetails.ProductID), orderItemDetails.Quantity)
		if err != nil {
			return models.OrderSuccessResponse{}, err
		}
	}
	return models.OrderSuccessResponse{
		OrderID:       OrderID,
		PaymentMethod: "Wallet",
		TotalAmount:   TotalAmount,
		PaymentStatus: "paid",
	}, nil
}

func (clean *OrderUseCase) ViewUserOrders(Token string) ([]models.ViewOrderDetails, error) {
	userId, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return []models.ViewOrderDetails{}, err
	}

	OrderDetails, err := clean.OrderRepo.GetOrderDetails(userId)
	if err != nil {
		return []models.ViewOrderDetails{}, err
	}
	return OrderDetails, nil
}

func (clean *OrderUseCase) CancelOrder(Token, orderId string, pid string) error {
	UserID, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return err
	}
	err = clean.OrderRepo.CheckOrder(orderId, UserID)
	if err != nil {
		return errors.New(`no orders found with this id`)
	}

	OrderDetails, err := clean.OrderRepo.CancelOrderDetails(UserID, orderId, pid)
	if err != nil {
		return err
	}

	if OrderDetails.OrderStatus == "Delivered" {
		return errors.New(`the order is delivered .Can't Cancel`)
	}

	if OrderDetails.OrderStatus == "Cancelled" {
		return errors.New(`the order is already cancelled`)
	}

	if OrderDetails.PaymentStatus == "paid" {
		err := clean.OrderRepo.ReturnAmountToWallet(UserID, orderId, pid)
		if err != nil {
			return err
		}

	}
	err = clean.OrderRepo.UpdateOrderFinalPrice(OrderDetails.OrderID, OrderDetails.TotalPrice)
	if err != nil {
		return err
	}
	proid, _ := strconv.Atoi(pid)
	err = clean.OrderRepo.UpdateStock(proid, OrderDetails.Quantity)
	if err != nil {
		return err
	}

	err = clean.OrderRepo.CancelOrder(orderId, pid, UserID)
	if err != nil {
		return err
	}

	return nil

}

func (clean *OrderUseCase) CancelOrderByAdmin(userID, order_id, pid int) error {
	orderID := strconv.Itoa(order_id)
	Pid := strconv.Itoa(pid)
	err := clean.OrderRepo.CheckOrder(orderID, uint(userID))
	fmt.Println(err)
	if err != nil {
		return errors.New(`no orders found with this id`)
	}
	err = clean.OrderRepo.CheckSingleOrder(Pid, orderID, uint(userID))
	if err != nil {
		return err
	}
	OrderDetails, err := clean.OrderRepo.CancelOrderDetails(uint(userID), orderID, Pid)
	if err != nil {
		return err
	}

	if OrderDetails.OrderStatus == "Delivered" {
		return errors.New(`the order is delivered .Can't Cancel`)
	}

	if OrderDetails.OrderStatus == "Cancelled" {
		return errors.New(`the order is already cancelled`)
	}

	if OrderDetails.PaymentStatus == "paid" {
		err := clean.OrderRepo.ReturnAmountToWallet(uint(userID), orderID, Pid)
		if err != nil {
			return err
		}

	}
	err = clean.OrderRepo.UpdateOrderFinalPrice(OrderDetails.OrderID, OrderDetails.TotalPrice)
	if err != nil {
		return err
	}

	err = clean.OrderRepo.UpdateStock(pid, OrderDetails.Quantity)
	if err != nil {
		return err
	}

	err = clean.OrderRepo.CancelOrder(orderID, Pid, uint(userID))
	if err != nil {
		return err
	}

	return nil
}

func (clean *OrderUseCase) ShipOrders(userID, orderId, pid int) error {
	orderID := strconv.Itoa(orderId)
	Pid := strconv.Itoa(pid)
	err := clean.OrderRepo.CheckOrder(orderID, uint(userID))
	fmt.Println(err)
	if err != nil {
		return errors.New(`no orders found with this id`)
	}
	err = clean.OrderRepo.CheckSingleOrder(Pid, orderID, uint(userID))
	if err != nil {
		return err
	}
	OrderStatus, err := clean.OrderRepo.GetOrderStatus(orderID, Pid)
	fmt.Println(OrderStatus)
	if err != nil {
		return err
	}
	if OrderStatus == "Cancelled" {
		return errors.New("the order is cancelled,cannot ship it")
	}

	if OrderStatus == "Delivered" {
		return errors.New("the order is already delivered")
	}

	if OrderStatus == "Shipped" {
		return errors.New("the order is already Shipped")
	}

	if OrderStatus == "pending" || OrderStatus == "processing" {
		err := clean.OrderRepo.ShipOrder(userID, orderId)
		if err != nil {
			return err
		}
		return nil
	}
	// if the shipment status is not processing or cancelled. Then it is defenetely cancelled
	return nil
}

func (clean *OrderUseCase) DeliverOrder(useriD, orderId, pid int) error {
	orderID := strconv.Itoa(orderId)
	Pid := strconv.Itoa(pid)
	err := clean.OrderRepo.CheckOrder(orderID, uint(useriD))
	fmt.Println(err)
	if err != nil {
		return errors.New(`no orders found with this id`)
	}
	err = clean.OrderRepo.CheckSingleOrder(Pid, orderID, uint(useriD))
	if err != nil {
		return err
	}
	OrderStatus, err := clean.OrderRepo.GetOrderStatus(orderID, Pid)
	if err != nil {
		return err
	}
	if OrderStatus == "Cancelled" {
		return errors.New("the order is cancelled,cannot deliver it")
	}

	if OrderStatus == "Delivered" {
		return errors.New("the order is already delivered")
	}

	if OrderStatus == "pending" {
		return errors.New("the order is not shipped yet")
	}

	if OrderStatus == "returned" {
		return errors.New(`the order is returned already by the customer`)
	}

	if OrderStatus == "Shipped" {
		err := clean.OrderRepo.DeliverOrder(useriD, orderID)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (clean *OrderUseCase) ReturnOrder(Token, orderID, pid string) error {
	UserID, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return err
	}

	err = clean.OrderRepo.CheckOrder(orderID, uint(UserID))
	fmt.Println(err)
	if err != nil {
		return errors.New(`no orders found with this id`)
	}

	err = clean.OrderRepo.CheckSingleOrder(pid, orderID, uint(UserID))
	if err != nil {
		return err
	}

	Order, err := clean.OrderRepo.GetOrderStatus(orderID, pid)
	if err != nil {
		return err
	}

	if Order != "returned" {
		return errors.New(`order is not delivered .Can't return`)
	}

	if Order == "returned" {
		return errors.New(`order is already returned`)
	}

	err = clean.OrderRepo.ReturnAmountToWallet(UserID, orderID, pid)
	if err != nil {
		return err
	}

	err = clean.OrderRepo.ReturnOrder(UserID, orderID, pid)
	if err != nil {
		return err
	}

	return nil
}

func (clean *OrderUseCase) ExecuteSalesReportByPeriod(period string) (*gofpdf.Fpdf, error) {
	startdate, enddate := utils.CalcualtePeriodDate(period)

	orders, err := clean.OrderRepo.GetByDate(startdate, enddate)
	if err != nil {
		return nil, errors.New("report fetching failed")
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Sales Report")
	pdf.Ln(10)
	pdf.Cell(0, 10, "Period:"+period)
	pdf.Ln(10)

	pdf.Cell(0, 10, "Total Sales: "+strconv.FormatFloat(orders.TotalSales, 'f', 2, 64))
	pdf.Ln(10)
	pdf.Cell(0, 10, "Total Orders: "+strconv.Itoa(int(orders.TotalOrders)))
	pdf.Ln(10)
	pdf.Cell(0, 10, "Average Order Price: "+strconv.FormatFloat(orders.AverageOrder, 'f', 2, 64))
	pdf.Ln(10)
	return pdf, nil
}

func (clean *OrderUseCase) ExecuteSalesReportByDate(startdate, enddate time.Time) (*gofpdf.Fpdf, error) {
	orders, err := clean.OrderRepo.GetByDate(startdate, enddate)
	if err != nil {
		return nil, errors.New("report fetching failed")
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Sales Report")
	pdf.Ln(10)
	pdf.Cell(0, 10, "From date: "+startdate.Format("02-01-2006"))
	pdf.Ln(10)
	pdf.Cell(0, 10, "To date: "+enddate.Format("02-01-2006"))
	pdf.Ln(10)

	pdf.Cell(0, 10, "Total Sales: "+strconv.FormatFloat(orders.TotalSales, 'f', 2, 64))
	pdf.Ln(10)
	pdf.Cell(0, 10, "Total Orders: "+strconv.Itoa(int(orders.TotalOrders)))
	pdf.Ln(10)
	pdf.Cell(0, 10, "Average Order Price: "+strconv.FormatFloat(orders.AverageOrder, 'f', 2, 64))
	pdf.Ln(10)
	return pdf, nil
}

func (clean *OrderUseCase) ExecuteSalesReportByPaymentMethod(startdate, enddate time.Time, paymentmethod string) (*gofpdf.Fpdf, error) {
	var payment string
	if paymentmethod == "1" {
		payment = "Cash On Delivery"
	} else if paymentmethod == "2" {
		payment = "Razorpay"
	} else if paymentmethod == "3" {
		payment = "Wallet"
	} else {
		return nil, errors.New(`payment method doesn't exist`)
	}
	orders, err := clean.OrderRepo.GetByPaymentMethod(startdate, enddate, paymentmethod)
	if err != nil {
		return nil, errors.New("report fetching failed")
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Sales Report")
	pdf.Ln(10)
	pdf.Cell(0, 10, "Payment Method: "+payment)
	pdf.Ln(10)
	pdf.Cell(0, 10, "From date: "+startdate.Format("02-01-2006"))
	pdf.Ln(10)
	pdf.Cell(0, 10, "To date: "+enddate.Format("02-01-2006"))
	pdf.Ln(10)

	pdf.Cell(0, 10, "Total Sales: "+strconv.FormatFloat(orders.TotalSales, 'f', 2, 64))
	pdf.Ln(10)
	pdf.Cell(0, 10, "Total Orders: "+strconv.Itoa(int(orders.TotalOrders)))
	pdf.Ln(10)
	pdf.Cell(0, 10, "Average Order Price: "+strconv.FormatFloat(orders.AverageOrder, 'f', 2, 64))
	pdf.Ln(10)

	return pdf, nil
}

func (clean *OrderUseCase) PrintInvoice(orderID int, Token string) (*gofpdf.Fpdf, error) {
	userID, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return nil, err
	}

	orde, err := clean.OrderRepo.GetOrderInvoice(orderID, int(userID))
	if err != nil {
		return nil, err
	}

	usr, err := clean.UserRepo.GetUserById(int(userID))
	if err != nil {
		return nil, err
	}

	usadres, err := clean.OrderRepo.GetAddressFromOrders(orde.AddressID, int(userID))
	if err != nil {
		return nil, err
	}

	items, err := clean.OrderRepo.GetProductDetailsFromOrders(orderID)
	if err != nil {
		return nil, err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Invoice")
	pdf.Ln(10)

	pdf.Cell(0, 10, "Customer Name: "+usr.Firstname)
	pdf.Ln(10)
	pdf.Cell(0, 10, "House Name: "+usadres.Housename)
	pdf.Ln(10)
	pdf.Cell(0, 10, "State: "+usadres.State)
	pdf.Ln(10)
	pdf.Cell(0, 10, "Phone: "+usadres.Phone)
	pdf.Ln(10)

	for _, item := range items {
		pdf.Cell(0, 10, "Item: "+item.Name)
		pdf.Ln(10)
		pdf.Cell(0, 10, "Price: "+strconv.FormatFloat(item.Price, 'f', 2, 64))
		pdf.Ln(10)
		pdf.Cell(0, 10, "Quantity: "+strconv.Itoa(item.Stock))
		pdf.Ln(10)
	}
	pdf.Ln(10)
	pdf.Cell(0, 10, "Total Amount: "+strconv.FormatFloat(float64(orde.TotalPrice), 'f', 2, 64))
	pdf.Ln(20)
	pdf.Cell(40, 10, "CityVibe: Thanks for shopping!")
	return pdf, nil
}

func (clean *OrderUseCase) ApplyCoupon(coupon, Token string) error {
	userID, err := utils.ExtractUserIdFromToken(Token)
	if err != nil {
		return err
	}

	err = clean.CouponRepo.CheckCouponUsage(userID, coupon)
	if err != nil {
		return err
	}
	DiscountRate, err := clean.CouponRepo.GetDiscountRate(coupon)
	if err != nil {
		return err
	}
	_, err = clean.OrderRepo.UpdateCartAmount(userID, uint(DiscountRate))
	if err != nil {
		return err
	}
	err = clean.CouponRepo.UpdateCouponUsage(userID, coupon)
	if err != nil {
		return err
	}
	err = clean.CouponRepo.UpdateCouponCount(coupon)
	if err != nil {
		return err
	}

	return nil
}

func (clean *OrderUseCase) SalesReportXL(start, end time.Time) (*excelize.File, error) {
	report, err := clean.OrderRepo.XLBYDATE(start, end)
	if err != nil {
		return nil, err
	}

	corereport, err := clean.OrderRepo.GetByDate(start, end)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheetName := "Sheet1"
	f.NewSheet(sheetName)

	f.SetColWidth("Sheet1", "A", "G", 20)
	// Set header style
	headerStyleID, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF", // Header text color
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4F81BD"}, // Header background color
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}
	styleID, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#DFEBF6"}, Pattern: 1},
	})
	if err != nil {
		return nil, err
	}
	f.SetCellValue(sheetName, "A1", fmt.Sprintf("CityVibe Sales Report (%s - %s)", start.Format("02-01-2006"), end.Format("02-01-2006")))
	f.SetCellStyle(sheetName, "A1", "A1", styleID)

	if err := f.MergeCell(sheetName, "A1", "G1"); err != nil {
		return nil, err
	}
	// Set header
	headers := []string{"Order number", "Customer Name", "Product Name", "Quantity", "Price"}
	for colIndex, header := range headers {
		cell := ConvertToAlphaString(colIndex+1) + "2"
		f.SetCellValue(sheetName, cell, header)

		// Apply header style
		f.SetCellStyle(sheetName, cell, cell, headerStyleID)

		// Auto adjust column width
		f.SetColWidth(sheetName, cell[:1], cell[:1], float64(len(header)*2)) // Adjust multiplier as needed
	}

	// Set data style
	dataStyleID, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	// Set data from report
	for rowIndex, record := range report {
		for colIndex, value := range []interface{}{record.OrderID, record.CustomerName, record.ProductName, record.Quantity, math.Round((float64(record.Price))*100) / 100} {
			cell := ConvertToAlphaString(colIndex+1) + fmt.Sprint(rowIndex+3)
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				return nil, err
			}

			// Apply data style
			f.SetCellStyle(sheetName, cell, cell, dataStyleID)

			// Auto adjust column width
			cellLetter := cell[:1]
			contentLength := len(fmt.Sprintf("%v", value)) * 2 // Adjust multiplier as needed
			currentWidth, _ := f.GetColWidth(sheetName, cellLetter)
			if contentLength > int(currentWidth) {
				f.SetColWidth(sheetName, cellLetter, cellLetter, float64(contentLength))
			}
		}
	}

	// Set total values style
	totalStyleID, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	// Set total values
	f.SetCellValue(sheetName, "F2", "Total Revenue Generated")
	f.SetCellValue(sheetName, "F3", "Total Orders")
	f.SetCellValue(sheetName, "F4", "Average Order Amount")
	f.SetCellValue(sheetName, "G2", corereport.TotalSales)
	f.SetCellValue(sheetName, "G3", corereport.TotalOrders)
	f.SetCellValue(sheetName, "G4", corereport.AverageOrder)
	f.SetCellValue(sheetName, "A1", fmt.Sprintf("CityVibe Sales Report (%s - %s)", start.Format("2006-01-02"), end.Format("2006-01-02")))

	// Apply total values style
	f.SetCellStyle(sheetName, "G2", "G4", totalStyleID)
	f.SetCellStyle(sheetName, "F2", "F4", totalStyleID)

	// Auto adjust column width for total values
	f.SetColWidth(sheetName, "G", "G", 15) // Adjust as needed
	f.SetColWidth(sheetName, "F", "F", 25) // Adjust as needed

	// ... rest of the code

	return f, nil
}

func  ConvertToAlphaString(index int) string {
	result := ""
	for index > 0 {
		index--
		result = string(rune('A'+index%26)) + result
		index /= 26
	}
	return result
}
