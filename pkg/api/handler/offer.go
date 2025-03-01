package handler

import (
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"ShowTimes/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OfferHandler struct {
	OfferUsecase interfaces.OfferUsecase
}

func NewOfferHandler(usecase interfaces.OfferUsecase) *OfferHandler {
	return &OfferHandler{
		OfferUsecase: usecase,
	}

}

// AddProductOffer adds a new product offer.
//
// @Summary Add Product Offer
// @Description Adds a new offer for a specific product.
// @Tags Admin Offer Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param productOffer body models.ProductOfferResp true "Product offer details in JSON format"
// @Success 201 {object} response.Response "Successfully added the product offer"
// @Failure 400 {object} response.Response "Invalid request format or constraints not satisfied"
// @Failure 500 {object} response.Response "Failed to add the product offer"
// @Router /admin/offer/product-offer [post]
func (oh *OfferHandler) AddProductOffer(c *gin.Context) {
	var productOffer models.ProductOfferResp

	if err := c.ShouldBindJSON(&productOffer); err != nil {
		errRes := response.ClientResponse(http.StatusBadGateway, "request fields are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	err := validator.New().Struct(productOffer)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "request fields are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err = oh.OfferUsecase.AddProductOffer(productOffer)

	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "could not add offer ", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusCreated, "Successfully added the offer", nil, nil)
	c.JSON(http.StatusCreated, successRes)
}

// AddCategoryOffer adds a new category-wide offer.
//
// @Summary Add Category Offer
// @Description Adds a new offer applicable to a specific category.
// @Tags Admin Offer Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param categoryOffer body models.CategorytOfferResp true "Category offer details in JSON format"
// @Success 201 {object} response.Response "Successfully added the category offer"
// @Failure 400 {object} response.Response "Invalid request format or constraints not satisfied"
// @Failure 500 {object} response.Response "Failed to add the category offer"
// @Router /admin/offer/category-offer [post]
func (of *OfferHandler) AddCategoryOffer(c *gin.Context) {

	var categoryOffer models.CategorytOfferResp

	if err := c.ShouldBindJSON(&categoryOffer); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "request fields in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	err := validator.New().Struct(categoryOffer)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err = of.OfferUsecase.AddCategoryOffer(categoryOffer)

	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "could not add offer in categorOffer", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Successfully added offer", nil, nil)
	c.JSON(http.StatusCreated, successRes)
}

// GetProductOffer retrieves all product offers.
//
// @Summary Get All Product Offers
// @Description Retrieves a list of all active product offers.
// @Tags Admin Offer Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.Response "Successfully retrieved all product offers"
// @Failure 400 {object} response.Response "Bad request: Unable to fetch product offers"
// @Router /admin/offer/product-offer [get]
func (of *OfferHandler) GetProductOffer(c *gin.Context) {

	products, err := of.OfferUsecase.GetProductOffer()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all offers", products, nil)
	c.JSON(http.StatusOK, successRes)
}

// GetCategoryOffer retrieves all category-wide offers.
//
// @Summary Get All Category Offers
// @Description Retrieves a list of all active category-wide offers.
// @Tags Admin Offer Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.Response "Successfully retrieved all category offers"
// @Failure 400 {object} response.Response "Bad request: Unable to fetch category offers"
// @Router /admin/offer/category-offer [get]
func (of *OfferHandler) GetCategoryOffer(c *gin.Context) {

	categories, err := of.OfferUsecase.GetCategoryOffer()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all offers", categories, nil)
	c.JSON(http.StatusOK, successRes)

}

// ExpireProductOffer invalidates a product-specific offer.
//
// @Summary Expire Product Offer
// @Description Marks a product offer as expired based on its ID.
// @Tags Admin Offer Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id query int true "Product offer ID to be expired"
// @Success 200 {object} response.Response "Successfully expired the product offer"
// @Failure 400 {object} response.Response "Bad request: Invalid offer ID format"
// @Failure 400 {object} response.Response "Bad request: Unable to expire the product offer"
// @Router /admin/offer/product-offer/expire [patch]
func (of *OfferHandler) ExpireProductOffer(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := of.OfferUsecase.ExpireProductOffer(id); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Coupon cannot be made invalid", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully made product offer invalid", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// ExpireCategoryOffer invalidates a category-wide offer.
//
// @Summary Expire Category Offer
// @Description Marks a category-wide offer as expired based on its ID.
// @Tags Admin Offer Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id query int true "Category offer ID"
// @Success 200 {object} response.Response "Successfully expired the category offer"
// @Failure 400 {object} response.Response "Bad request: Invalid offer ID format"
// @Failure 400 {object} response.Response "Bad request: Unable to expire the category offer"
// @Router /admin/offer/category-offer/expire [patch]
func (of *OfferHandler) ExpireCategoryOffer(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := of.OfferUsecase.ExpireCategoryOffer(id); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Coupon cannot be made invalid", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully made category offer invalid", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
