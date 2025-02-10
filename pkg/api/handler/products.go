package handler

import (
	"ShowTimes/pkg/domain"
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"ShowTimes/pkg/utils/response"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	ProductUseCase interfaces.ProductUseCase
}

func NewProductHandler(usecase interfaces.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		ProductUseCase: usecase,
	}

}

// AddProduct adds a new product.
// @Summary Add product
// @Description Adds a new product using the provided details and image.
// @Tags Product
// @Accept multipart/form-data
// @Produce json
// @Security BearerTokenAuth
// @Param category_id formData integer true "Category ID"
// @Param product_name formData string true "Product name"
// @Param color formData string true "Product color"
// @Param stock formData integer true "Product stock"
// @Param price formData number true "Product price"
// @Param image formData file true "Product image"
// @Success 200 {object} response.Response  "Success: Product added successfully"
// @Failure 400 {object} response.Response  "Bad request: Retrieving image error or could not add the product"
// @Router /admin/product [post]

func (i *ProductHandler) AddProducts(c *gin.Context) {
	var products models.AddProducts
	cat := c.PostForm("category_id")
	products.CategoryID, _ = strconv.Atoi(cat)
	products.ProductName = c.PostForm("product_name")
	products.Color = c.PostForm("color")
	products.Stock, _ = strconv.Atoi(c.PostForm("stock"))
	products.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)

	file, err := c.FormFile("image")
	if err != nil {
		errorResp := response.ClientResponse(http.StatusBadRequest, "retrieving image from the Form error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}
	ProductResponse, err := i.ProductUseCase.AddProducts(products, file)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "couldnot add the ", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfully added the product", ProductResponse, nil)
	c.JSON(http.StatusOK, successResp)
}

// ListProducts lists products with pagination.
// @Summary List products
// @Description Retrieves a paginated list of products.
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param page query integer false "Page number (default: 1)"
// @Param per_page query integer false "Number of products per page (default: 5)"
// @Success 200 {object} response.Response  "Success: Products displayed successfully"
// @Failure 400 {object} response.Response  "Bad request: Product display error"
// @Router /admin/product [get]

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
	product_list, err := i.ProductUseCase.ListProducts(pageNoInt, pageListInt)

	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Product cannot be displayed...", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
	}
	message := "Product list"

	successResp := response.ClientResponse(http.StatusOK, message, product_list, nil)
	c.JSON(http.StatusOK, successResp)

}

// EditProduct updates an existing product.
// @Summary Edit product
// @Description Updates an existing product using the provided details.
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param Product body domain.Product true "Product details to be updated"
// @Success 200 {object} response.Response  "Success: Product edited successfully"
// @Failure 400 {object} response.Response  "Bad request: Fields are in the wrong format or could not edit the product"
// @Router /admin/product [patch]

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
	modInventory, err := u.ProductUseCase.EditProducts(inventory, idInt)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "couldnot edit the product", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfully edited", modInventory, nil)
	c.JSON(http.StatusOK, successResp)
}

// DeleteProduct deletes an existing product by ID.
// @Summary Delete product
// @Description Deletes an existing product by the provided ID.
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id query string true "Product ID to be deleted"
// @Success 200 {object} response.Response  "Success: Product deleted successfully"
// @Failure 400 {object} response.Response  "Bad request: Product ID provided in wrong format or deletion error"
// @Router /admin/product [delete]

func (u *ProductHandler) DeleteProducts(c *gin.Context) {
	inventoryID := c.Query("id")

	err := u.ProductUseCase.DeleteProducts(inventoryID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully deleted theproduct", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// UpdateProduct updates the stock of an existing product.
// @Summary Update product stock
// @Description Updates the stock of an existing product using the provided details.
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param ProductUpdate body models.ProductUpdate true "Product details for stock update"
// @Success 200 {object} response.Response  "Success: Product stock updated successfully"
// @Failure 400 {object} response.Response  "Bad request: Fields provided in wrong format or could not update the product stock"
// @Router /product/stock [patch]

func (i *ProductHandler) UpdateProducts(c *gin.Context) {
	var p models.ProductUpdate

	if err := c.BindJSON(&p); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	a, err := i.ProductUseCase.UpdateProducts(p.Productid, p.Stock)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Coud not update the inventory stock", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	succesResp := response.ClientResponse(http.StatusOK, "Successfully updated inventory stock", a, nil)
	c.JSON(http.StatusOK, succesResp)

}
