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
