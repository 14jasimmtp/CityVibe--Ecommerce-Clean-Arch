package usecase

import (
	"errors"
	"testing"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	mock "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repositorymock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCategoryRepo := mock.NewMockCategoryRepo(ctrl)
	categoryUseCase :=NewCategoryUsecase(mockCategoryRepo)

	expectedCategories := []domain.Category{
		{ID: 1, Category: "Category1"},
		{ID: 2, Category: "Category2"},
	}

	t.Run("Success", func(t *testing.T) {
		mockCategoryRepo.EXPECT().GetCategory().Return(expectedCategories, nil)

		result, err := categoryUseCase.GetCategory()
		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockError := errors.New("mock error")
		mockCategoryRepo.EXPECT().GetCategory().Return([]domain.Category{}, mockError)

		result, err := categoryUseCase.GetCategory()

		assert.Error(t, err)
		assert.Equal(t, []domain.Category{}, result)
		assert.EqualError(t, err, mockError.Error())
	})
}


