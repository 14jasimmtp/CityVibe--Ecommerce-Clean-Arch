package interfaceRepo

import (
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
)

type CategoryRepo interface {
	GetCategory() ([]domain.Category, error)
	AddCategory(category models.Category) (domain.Category, error)
	DeleteCategory(id string) error
	UpdateCategory(current string, new string) (domain.Category, error)
	CheckCategory(current string) (bool, error)
}
