package interfaceUsecase

import (
	"mime/multipart"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
)

type ProductUsecase interface {
	AddProduct(product models.AddProduct, image *multipart.FileHeader) (models.UpdateProduct, error)
	GetAllProducts() ([]models.Product, error)
	EditProductDetails(id string, product models.AddProduct) (models.UpdateProduct, error)
	DeleteProduct(id string) error
	ShowProductsByCategory() ([]domain.Product, error)
	SeeAllProducts() ([]domain.Product, error)
	GetSingleProduct(id string) (models.Product, error)
	FilterProductCategoryWise(category string) ([]models.Product, error)
	SearchProduct(search string) ([]models.Product, error)
	FilterProducts(category, size string, minPrice, maxPrice float64) ([]models.UpdateProduct, error)
}
