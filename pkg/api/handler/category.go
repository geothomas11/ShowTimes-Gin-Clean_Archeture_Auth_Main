package handler

import (
	"ShowTimes/pkg/domain"
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"ShowTimes/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	CategoryUseCase interfaces.CategoryUseCase
}

func NewCategoryHandler(usecase interfaces.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		CategoryUseCase: usecase,
	}
}

// AddCategory adds a new category.
// @Summary Add a new category
// @Description Adds a new category based on the provided details.
// @Tags Category Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param AddCategory body domain.Category true "Category details to add"
// @Success 200 {object} response.Response "Success: Category added successfully"
// @Failure 400 {object} response.Response "Bad request: Fields are provided in the wrong format"
// @Failure 401 {object} response.Response "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response "Internal server error: Could not add the category"
// @Router /admin/category [post]
func (cat *CategoryHandler) AddCategory(c *gin.Context) {
	var category domain.Category
	if err := c.BindJSON(&category); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Fields are provided in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	CategoryResponse, err := cat.CategoryUseCase.AddCategory(category)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Could not add the category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Category added successfully", CategoryResponse, nil)
	c.JSON(http.StatusOK, successRes)
}

// GetCategory retrieves all categories.
// @Summary Get all categories
// @Description Retrieves a list of all available categories.
// @Tags Category Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.Response "Success: Retrieved categories successfully"
// @Failure 400 {object} response.Response "Bad request: Fields provided are in the wrong format"
// @Failure 401 {object} response.Response "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response "Internal server error: Could not retrieve categories"
// @Router /admin/categories [get]
func (cat *CategoryHandler) GetCategory(c *gin.Context) {
	categories, err := cat.CategoryUseCase.GetCategories()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully retrieved all categories", categories, nil)
	c.JSON(http.StatusOK, successRes)
}

// UpdateCategory updates the name of an existing category.
// @Summary Update category name
// @Description Updates the name of an existing category based on the provided details.
// @Tags Category Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param UpdateCategory body models.SetNewName true "Current and new category name"
// @Success 200 {object} response.Response "Success: Category updated successfully"
// @Failure 400 {object} response.Response "Bad request: Fields provided are in the wrong format"
// @Failure 401 {object} response.Response "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response "Internal server error: Could not update the category"
// @Router /admin/category/update [patch]
func (cat *CategoryHandler) UpdateCategory(c *gin.Context) {
	var p models.SetNewName

	if err := c.BindJSON(&p); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	updatedCategory, err := cat.CategoryUseCase.UpdateCategory(p.Current, p.New)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not update the category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Category updated successfully", updatedCategory, nil)
	c.JSON(http.StatusOK, successRes)
}

// DeleteCategory deletes a category by ID.
// @Summary Delete a category
// @Description Deletes a category based on the provided category ID.
// @Tags Category Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id query string true "Category ID to delete"
// @Success 200 {object} response.Response "Success: Category deleted successfully"
// @Failure 400 {object} response.Response "Bad request: Fields are not provided in the correct format"
// @Failure 401 {object} response.Response "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response "Internal server error: Could not delete the category"
// @Router /admin/category/delete [delete]
func (cat *CategoryHandler) DeleteCategory(c *gin.Context) {
	categoryID := c.Query("id")
	err := cat.CategoryUseCase.DeleteCategory(categoryID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Category deleted successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
