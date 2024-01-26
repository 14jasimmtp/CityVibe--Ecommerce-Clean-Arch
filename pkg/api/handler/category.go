package handlers

import (
	"net/http"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	interfaceUsecase "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/usecase/interface"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	CategoryUsecase interfaceUsecase.CategoryUsecase
}

func NewCategoryHandler(usecase interfaceUsecase.CategoryUsecase) *CategoryHandler{
	return &CategoryHandler{CategoryUsecase: usecase}
}

// AddCategory godoc
// @Summary Add a new category
// @Description Add a new category using the provided details.
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Param category body models.Category true "Details of the category to be added"
// @Success 200 {object} string "message": "Successfully added category", "category": Cate
// @Failure 400 {object} string "error": "Bad Request"
// @Failure 500 {object} string "error": "Internal Server Error"
// @Router /admin/category [post]
func (clean *CategoryHandler) AddCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter details correctly"})
		return
	}
	Cate, err := clean.CategoryUsecase.AddCategory(category)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully added category", "category": Cate})

}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category based on the provided category ID.
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Param id query string true "Category ID to be deleted"
// @Success 200 {object} string "message": "Successfully deleted category"
// @Failure 500 {object} string "error": "Internal Server Error"
// @Router /admin/category [delete]
func (clean *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Query("id")
	err := clean.CategoryUsecase.DeleteCategory(id)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted category"})
}

// GetCategory godoc
// @Summary Get all categories
// @Description Retrieve all categories.
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Success 200 {object} string "message": "Categories", "categories": category
// @Failure 400 {object} string "error": "Bad Request"
// @Failure 500 {object} string "error": "Internal Server Error"
// @Router /admin/category [get]
func (clean *CategoryHandler) GetCategory(c *gin.Context) {
	category, err := clean.CategoryUsecase.GetCategory()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "categories", "categories": category})
}

// UpdateCategory godoc
// @Summary Update category name
// @Description Update the name of a category based on the provided details.
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Param categoryUpdate body models.SetNewName true "Current and new names for the category"
// @Success 200 {object} string "message": "Category updated successfully", "updated category": Newcate
// @Failure 400 {object} string "error": "Bad Request"
// @Failure 500 {object} string "error": "Internal Server Error"
// @Router /admin/category [put]
func (clean *CategoryHandler) UpdateCategory(c *gin.Context) {
	var categoryUpdate models.SetNewName
	if err := c.ShouldBindJSON(&categoryUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Newcate, err := clean.CategoryUsecase.UpdateCategory(categoryUpdate.Current, categoryUpdate.New)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "category updated successfully", "updated category": Newcate})
}
