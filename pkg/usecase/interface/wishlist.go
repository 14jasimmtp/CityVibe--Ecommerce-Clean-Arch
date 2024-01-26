package interfaceUsecase

import "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"

type WishlistUsecase interface {
	AddProductToWishlist(pid string, Token string) error
	ViewUserWishlist(Token string) ([]models.UpdateProduct, error)
	RemoveProductFromWishlist(pid, Token string) error
}
