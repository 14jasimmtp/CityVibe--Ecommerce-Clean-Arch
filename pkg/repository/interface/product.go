package interfaceRepo

import (
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
)

type ProductRepo interface {
	AddProduct(product models.AddProduct) (models.UpdateProduct, error)
	EditProductDetails(id string, product models.AddProduct) (models.UpdateProduct, error)
	DeleteProduct(id int) error
	GetAllProducts() ([]models.Product, error)
	SeeAllProducts() ([]domain.Product, error)
	GetSingleProduct(id string) (models.Product, error)
	FilterProductCategoryWise(category string) ([]models.Product, error)
	CheckStock(pid int) error
	GetProductAmountFromID(pid string) (float64, error)
	SearchProduct(search string) ([]models.Product, error)
	FilterProducts(category, size string, minPrice, maxPrice float64) ([]models.UpdateProduct, error)
	GetAllOffers() ([]models.Offer, error)
	GetProductsByCategoryoffer(id int) ([]models.Product, error)
	GetProductById(id int) (*models.Product, error)
	UpdateProduct(product *models.Product) error
}
