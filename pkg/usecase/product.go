package usecase

import (
	"fmt"
	"mime/multipart"
	"strconv"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	interfaceRepo "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repository/interface"
	interfaceUsecase "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/usecase/interface"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/utils"
)

type ProductUseCase struct {
	ProductRepo interfaceRepo.ProductRepo
}

func NewProductUsecase(repo interfaceRepo.ProductRepo) interfaceUsecase.ProductUsecase{
	return &ProductUseCase{ProductRepo: repo}
}

func (clean *ProductUseCase) AddProduct(product models.AddProduct, image *multipart.FileHeader) (models.UpdateProduct, error) {
	sess := utils.CreateSession()
	// fmt.Println("sess", sess)
	
	ImageURL, err := utils.UploadImageToS3(image, sess)
	if err != nil {
		fmt.Println("err:", err)
		return models.UpdateProduct{}, err
	}
	fmt.Println("err:", err)
	product.ImageURL = ImageURL
	fmt.Println("image,", ImageURL)

	ProductResponse, err := clean.ProductRepo.AddProduct(product)
	if err != nil {
		return models.UpdateProduct{}, err
	}
	return ProductResponse, nil
}

func (clean *ProductUseCase) GetAllProducts() ([]models.Product, error) {
	ProductDetails, err := clean.ProductRepo.GetAllProducts()
	if err != nil {
		return []models.Product{}, err
	}
	return ProductDetails, nil
}

func (clean *ProductUseCase) EditProductDetails(id string, product models.AddProduct) (models.UpdateProduct, error) {
	UpdatedProduct, err := clean.ProductRepo.EditProductDetails(id, product)
	if err != nil {
		return models.UpdateProduct{}, err
	}
	return UpdatedProduct, nil
}

func (clean *ProductUseCase) DeleteProduct(id string) error {
	idnum, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	err = clean.ProductRepo.DeleteProduct(idnum)
	if err != nil {
		return err
	}

	return nil
}

func (clean *ProductUseCase) ShowProductsByCategory() ([]domain.Product, error) {

	return []domain.Product{}, nil
}

func (clean *ProductUseCase) SeeAllProducts() ([]domain.Product, error) {
	products, err := clean.ProductRepo.SeeAllProducts()
	if err != nil {
		return []domain.Product{}, err
	}
	return products, nil
}

func (clean *ProductUseCase) GetSingleProduct(id string) (models.Product, error) {
	product, err := clean.ProductRepo.GetSingleProduct(id)
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (clean *ProductUseCase) FilterProductCategoryWise(category string) ([]models.Product, error) {
	products, err := clean.ProductRepo.FilterProductCategoryWise(category)
	if err != nil {
		return []models.Product{}, err
	}
	return products, nil
}

func (clean *ProductUseCase) SearchProduct(search string) ([]models.Product, error) {
	products, err := clean.ProductRepo.SearchProduct(search)
	if err != nil {
		return []models.Product{}, err
	}

	return products, nil
}

func (clean *ProductUseCase) FilterProducts(category, size string, minPrice, maxPrice float64) ([]models.UpdateProduct, error) {
	filteredProducts, err := clean.ProductRepo.FilterProducts(category, size, minPrice, maxPrice)
	if err != nil {
		return []models.UpdateProduct{}, err
	}
	return filteredProducts, nil
}

