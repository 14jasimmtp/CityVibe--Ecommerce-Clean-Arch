package interfaceRepo

import "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"

type CartRepo interface {
	AddToCart(pid int, userid uint, productAmount float64) error
	DisplayCart(userid uint) ([]models.Cart, error)
	RemoveProductFromCart(pid int, userid uint) error
	CheckProductExistInCart(userId uint, pid string) (bool, error)
	UpdateQuantity(userid uint, pid, quantity string) ([]models.Cart, error)
	CartTotalAmount(userid uint) (float64, error)
	CartFinalPrice(userid uint) (float64, error)
	CheckCartExist(userID uint) bool
	EmptyCart(userID uint) error
	ProductQuantityCart(userID uint, pid string) (int, error)
	UpdateQuantityAdd(id uint, prdt_id string) error
	UpdateQuantityless(id uint, prdt_id string) error
	CartExist(userID uint) (bool, error)
	UpdateCart(quantity int, price float64, userID uint, product_id string) error
	UpdateTotalPrice(id uint, product_id string) error
	TotalPrizeOfProductInCart(userID uint, productID string) (float64, error)
	CheckSingleProduct(id string) (models.Product, error)
	CheckCartStock(pid int) error
	GetCartProductAmountFromID(pid string) (float64, error)
}
