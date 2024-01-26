package usecase

import (
	"errors"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	interfaceRepo "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repository/interface"
	interfaceUsecase "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/usecase/interface"
)

type CategoryUseCase struct {
	CategoryRepo interfaceRepo.CategoryRepo
}

func NewCategoryUsecase(repo interfaceRepo.CategoryRepo) interfaceUsecase.CategoryUsecase{
	return &CategoryUseCase{CategoryRepo: repo}
}

func (clean *CategoryUseCase) GetCategory() ([]domain.Category, error) {
	category, err := clean.CategoryRepo.GetCategory()
	if err != nil {
		return []domain.Category{}, err
	}
	return category, nil

}
func (clean *CategoryUseCase) AddCategory(category models.Category) (domain.Category, error) {
	categories, err := clean.CategoryRepo.AddCategory(category)
	if err != nil {
		return domain.Category{}, err
	}
	return categories, nil
}
func (clean *CategoryUseCase) UpdateCategory(current string, new string) (domain.Category, error) {
	categries, err := clean.CategoryRepo.CheckCategory(current)
	if err != nil {
		return domain.Category{}, err
	}
	if !categries {
		return domain.Category{}, errors.New("category doesn't exist")
	}
	newCate, err := clean.CategoryRepo.UpdateCategory(current, new)
	if err != nil {
		return domain.Category{}, err
	}
	return newCate, nil
}
func (clean *CategoryUseCase) DeleteCategory(id string) error {
	err := clean.CategoryRepo.DeleteCategory(id)
	if err != nil {
		return err
	}
	return nil
}
