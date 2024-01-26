package interfaceUsecase

import "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"

type CartUsecase interface {
	ViewCart(Token string) (models.CartResponse, error)
	AddToCart(pid, Token string) (models.CartResponse, error)
	RemoveProductsFromCart(pid, Token string) (models.CartResponse, error)
	UpdateQuantityFromCart(Token, pid, quantity string) ([]models.Cart, error)
	UpdateQuantityIncrease(Token, pid string) error
	UpdatePriceAdd(Token, pid string) error
	UpdateQuantityDecrease(Token, pid string) error
	UpdatePriceDecrease(Token, pid string) error
	EraseCart(Token string) (models.CartResponse, error)
}
