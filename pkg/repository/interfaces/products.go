package interfaces

import (
	"ShowTimes/pkg/utils/models"
)

type ProductRepository interface {
	AddProducts(inventory models.AddProducts, url string) (models.ProductResponse, error)
	ListProducts(int, int) ([]models.ProductUserResponse, error)
	// EditProducts(domain.Product, int) (domain.Product, error)
	DeleteProducts(id string) error
	CheckProducts(p_id int) (bool, error)
	UpdateProducts(p_id int, stock int) (models.ProductResponse, error)
	CheckProductAvailable(product_id int) (bool, error)
	GetPriceOfProduct(product_id int) (float64, error)
	EditProduct(product models.ProductEdit) (models.ProductUserResponse, error)
}
