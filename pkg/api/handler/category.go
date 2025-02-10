package handler

import (
	// "ShowTimes/pkg/usecase"
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
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param AddCategory body domain.Category true "Category details to add"
// @Success 200 {object} response.Response  "Success: Category added successfully"
// @Failure 400 {object} response.Response  "Bad request: Fields are provided in the wrong format"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Could not add the category"
// @Router /admin/category [post]

func (cat *CategoryHandler) AddCategory(c *gin.Context) {
	var category domain.Category
	if err := c.BindJSON(&category); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields are provided in wrong format ", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	CatgoryResponse, err := cat.CategoryUseCase.AddCategory(category)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Coudnot not add the category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Category added successfully", CatgoryResponse, nil)
	c.JSON(http.StatusOK, successRes)

}
func (Cat *CategoryHandler) GetCategory(c *gin.Context) {

	categories, err := Cat.CategoryUseCase.GetCategories()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all categories", categories, nil)
	c.JSON(http.StatusOK, successRes)

}

// GetCategory retrieves all categories.
// @Summary Retrieve all categories
// @Description Retrieves all categories available.
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.Response  "Success: Retrieved all categories successfully"
// @Failure 400 {object} response.Response  "Bad request: Fields provided in the wrong format"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Could not retrieve categories"
// @Router /admin/category [get]

func (cat *CategoryHandler) UpdateCategory(c *gin.Context) {
	var p models.SetNewName

	if err := c.BindJSON(&p); err != nil {

		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	up_category, err := cat.CategoryUseCase.UpdateCategory(p.Current, p.New)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Cannot update ", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Categories updated successfully!!!", up_category, nil)
	c.JSON(http.StatusOK, successRes)

}

// DeleteCategory deletes a category by ID.
// @Summary Delete category
// @Description Deletes a category based on the provided category ID.
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id query string true "Category ID to delete"
// @Success 200 {object} YourResponseObject "Success: Category deleted successfully"
// @Failure 400 {object} YourResponseObject "Bad request: Fields are not provided in the correct format"
// @Failure 401 {object} YourResponseObject "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} YourResponseObject "Internal server error: Could not delete the category"
// @Router /categories/delete [delete]

func (cat *CategoryHandler) DeleteCategory(c *gin.Context) {
	categoryID := c.Query("id")
	err := cat.CategoryUseCase.DeleteCategory(categoryID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields  provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully Deleted...", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
