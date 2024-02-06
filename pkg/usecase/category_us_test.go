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

	// Define your test data
	expectedCategories := []domain.Category{
		{ID: 1, Category: "Category1"},
		{ID: 2, Category: "Category2"},
	}

	// Test case for success scenario
	t.Run("Success", func(t *testing.T) {
		// Set up expectations for your mock CategoryRepo
		mockCategoryRepo.EXPECT().GetCategory().Return(expectedCategories, nil)

		// Call the method being tested
		result, err := categoryUseCase.GetCategory()

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, result)
	})

	// Test case for error scenario
	t.Run("Error", func(t *testing.T) {
		// Set up expectations for your mock CategoryRepo
		mockError := errors.New("mock error")
		mockCategoryRepo.EXPECT().GetCategory().Return([]domain.Category{}, mockError)

		// Call the method being tested
		result, err := categoryUseCase.GetCategory()

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, []domain.Category{}, result)
		assert.EqualError(t, err, mockError.Error())
	})
}
