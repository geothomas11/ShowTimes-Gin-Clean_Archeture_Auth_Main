package interfaces

import (
	"ShowTimes/pkg/domain"
	"ShowTimes/pkg/utils/models"
	"mime/multipart"
)

type ProductUseCase interface {
	AddProducts(inventory models.AddProducts,file*multipart.FileHeader) (models.ProductResponse, error)
	ListProducts(int, int) ([]models.ProductUserResponse, error)
	EditProducts(domain.Product, int) (domain.Product, error)
	DeleteProducts(id string) error
	UpdateProducts(poductID int, stock int) (models.ProductResponse, error)
}
