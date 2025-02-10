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
//
// @Summary Add a new product
// @Description Adds a new product using the provided details and image.
// @Tags Product
// @Accept multipart/form-data
// @Produce json
// @Security BearerTokenAuth
// @Param category_id formData integer true "Category ID of the product"
// @Param product_name formData string true "Name of the product"
// @Param color formData string true "Color of the product"
// @Param stock formData integer true "Stock quantity of the product"
// @Param price formData number true "Price of the product"
// @Param image formData file true "Image file of the product"
// @Success 200 {object} response.Response "Product added successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid input or image upload error"
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
		errorResp := response.ClientResponse(http.StatusBadRequest, "Error retrieving image from the form", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}
	ProductResponse, err := i.ProductUseCase.AddProducts(products, file)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Could not add the product", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfully added the product", ProductResponse, nil)
	c.JSON(http.StatusOK, successResp)
}

// ListProducts retrieves a paginated list of products.
//
// @Summary List products
// @Description Fetches a paginated list of products.
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param page query integer false "Page number (default: 1)"
// @Param per_page query integer false "Number of products per page (default: 5)"
// @Success 200 {object} response.Response "Product list retrieved successfully"
// @Failure 400 {object} response.Response "Bad request: Error displaying products"
// @Router /admin/product [get]
func (i *ProductHandler) ListProducts(c *gin.Context) {
	pageNo := c.DefaultQuery("page", "1")
	pageList := c.DefaultQuery("per_page", "5")
	pageNoInt, err := strconv.Atoi(pageNo)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Invalid page number format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	pageListInt, err := strconv.Atoi(pageList)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Invalid per_page format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	productList, err := i.ProductUseCase.ListProducts(pageNoInt, pageListInt)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Error retrieving product list", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Product list retrieved successfully", productList, nil)
	c.JSON(http.StatusOK, successResp)
}

// EditProduct updates the details of an existing product.
//
// @Summary Edit product
// @Description Updates an existing product using the provided details.
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param inventory_id query integer true "Inventory ID of the product to update"
// @Param product body domain.Product true "Updated product details"
// @Success 200 {object} response.Response "Product updated successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid input or product update error"
// @Router /admin/product [patch]
func (u *ProductHandler) EditProducts(c *gin.Context) {
	var inventory domain.Product
	id := c.Query("inventory_id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Invalid inventory ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	if err := c.BindJSON(&inventory); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Invalid input format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	modInventory, err := u.ProductUseCase.EditProducts(inventory, idInt)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Error updating the product", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Product updated successfully", modInventory, nil)
	c.JSON(http.StatusOK, successResp)
}

// DeleteProduct removes a product by its ID.
//
// @Summary Delete product
// @Description Deletes an existing product using its ID.
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id query integer true "ID of the product to be deleted"
// @Success 200 {object} response.Response "Product deleted successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid product ID or deletion error"
// @Router /admin/product [delete]
func (u *ProductHandler) DeleteProducts(c *gin.Context) {
	inventoryID := c.Query("id")

	err := u.ProductUseCase.DeleteProducts(inventoryID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid product ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Product deleted successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// UpdateProduct updates the stock quantity of a product.
//
// @Summary Update product stock
// @Description Updates the stock quantity of an existing product.
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param request body models.ProductUpdate true "Updated stock details"
// @Success 200 {object} response.Response "Product stock updated successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid input or update error"
// @Router /product/stock [patch]
func (i *ProductHandler) UpdateProducts(c *gin.Context) {
	var p models.ProductUpdate

	if err := c.BindJSON(&p); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid input format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	updatedProduct, err := i.ProductUseCase.UpdateProducts(p.Productid, p.Stock)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Error updating inventory stock", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Inventory stock updated successfully", updatedProduct, nil)
	c.JSON(http.StatusOK, successResp)
}
