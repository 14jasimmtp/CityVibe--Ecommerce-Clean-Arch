package interfaceUsecase

import (
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
)

type CategoryUsecase interface {
	GetCategory() ([]domain.Category, error)
	AddCategory(category models.Category) (domain.Category, error)
	UpdateCategory(current string, new string) (domain.Category, error)
	DeleteCategory(id string) error
}
