package handler

import (
	"ShowTimes/pkg/domain"
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"ShowTimes/pkg/utils/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	InventoryUseCase interfaces.ProductUseCase
}

func NewProductHandler(usecase interfaces.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		InventoryUseCase: usecase,
	}

}

func (i *ProductHandler) AddProducts(c *gin.Context) {
	var inventory models.AddProducts

	if err := c.ShouldBindJSON(&inventory); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "form file error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		fmt.Println("error", err)
		return
	}
	InventoryResponse, err := i.InventoryUseCase.AddProducts(inventory)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "could not add the inventory", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfilly invetory added ", InventoryResponse, nil)
	c.JSON(http.StatusOK, successResp)

}

func (i *ProductHandler) ListProducts(c *gin.Context) {
	pageNo := c.DefaultQuery("page", "1")
	pageList := c.DefaultQuery("per_page", "5")
	pageNoInt, err := strconv.Atoi(pageNo)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Product Cannot be displayed.", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	pageListInt, err := strconv.Atoi(pageList)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Product cannot be displayed..", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
	}
	product_list, err := i.InventoryUseCase.ListProducts(pageNoInt, pageListInt)

	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Product cannot be displayed...", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
	}
	message := "Product list"

	successResp := response.ClientResponse(http.StatusOK, message, product_list, nil)
	c.JSON(http.StatusOK, successResp)

}

func (u *ProductHandler) EditProducts(c *gin.Context) {
	var inventory domain.Product

	id := c.Query("inventory_id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Problems with the id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	if err := c.BindJSON(&inventory); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields are in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	modInventory, err := u.InventoryUseCase.EditProducts(inventory, idInt)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "couldnot edit the product", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfully edited", modInventory, nil)
	c.JSON(http.StatusOK, successResp)
}

func (u *ProductHandler) DeleteProducts(c *gin.Context) {
	inventoryID := c.Query("id")

	err := u.InventoryUseCase.DeleteProducts(inventoryID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully deleted theproduct", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (i *ProductHandler) UpdateProducts(c *gin.Context) {
	var p models.ProductUpdate

	if err := c.BindJSON(&p); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	a, err := i.InventoryUseCase.UpdateProducts(p.Productid, p.Stock)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Coud not update the inventory stock", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	succesResp := response.ClientResponse(http.StatusOK, "Successfully updated inventory stock", a, nil)
	c.JSON(http.StatusOK, succesResp)

}
