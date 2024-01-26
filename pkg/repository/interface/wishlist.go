package interfaceRepo

import "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"

type WishlistRepo interface {
	CheckExistInWishlist(userID uint, pid string) error
	AddProductToWishlist(pid string, userID uint) error
	GetWishlistProducts(userID uint) ([]models.UpdateProduct, error)
	RemoveProductFromWishlist(pid string, userID uint) error
	WishlistSingleProduct(id string) (models.Product, error)
}
