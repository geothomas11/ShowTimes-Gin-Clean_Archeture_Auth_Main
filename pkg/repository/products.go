package repository

import (
	interfaces "ShowTimes/pkg/repository/interfaces"
	"ShowTimes/pkg/utils/errmsg"
	"ShowTimes/pkg/utils/models"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) interfaces.ProductRepository {
	return &ProductRepository{
		DB: db,
	}

}
func (i *ProductRepository) AddProducts(product models.AddProducts, url string) (models.ProductResponse, error) {

	var count int64
	i.DB.Model(&models.Product{}).Where("product_name = ? AND category_id = ?", product.ProductName, product.CategoryID).Count(&count)
	if count > 0 {

		return models.ProductResponse{}, errors.New(errmsg.ErrProductExistTrue)
	}

	if product.Stock < 0 || product.Price < 0 {
		return models.ProductResponse{}, errors.New("stock and price" + errmsg.ErrDataNegative)
	}

	query := `
        INSERT INTO products (category_id, product_name, color, stock, price)
        VALUES (?, ?, ?, ?, ?);
    `
	err := i.DB.Exec(query, product.CategoryID, product.ProductName, product.Color, product.Stock, product.Price).Error
	if err != nil {
		return models.ProductResponse{}, err
	}

	// getting inserted product detailsS
	var productResponse models.ProductResponse

	query = "SELECT id,category_id,product_name,color,stock,price FROM products where  product_name = ? AND category_id = ?"
	errr := i.DB.Raw(query, product.ProductName, product.CategoryID).Scan(&productResponse).Error

	if errr != nil {
		return productResponse, errors.New(errmsg.ErrGetDB)
	}
	//Adding url to image table

	queryimage := "INSERT INTO product_images (product_id, url) VALUES (?, ?)"

	imgErr := i.DB.Exec(queryimage, productResponse.ID, url).Error
	if imgErr != nil {

		return models.ProductResponse{}, imgErr
	}
	return productResponse, nil
}

func (prod *ProductRepository) ListProducts(pageList, offset int) ([]models.ProductUserResponse, error) {
	var product_list []models.ProductUserResponse

	query := "SELECT i.id,i.category_id,c.category,i.product_name,i.color,i.price,i.url FROM products AS i INNER JOIN categories c ON i.category_id = c.id LIMIT $1 OFFSET $2"
	err := prod.DB.Raw(query, pageList, offset).Scan(&product_list).Error

	if err != nil {
		return []models.ProductUserResponse{}, errors.New("error checking the product details")
	}
	return product_list, nil

}

func (db *ProductRepository) EditProduct(product models.ProductEdit) (models.ProductUserResponse, error) {
	var modProduct models.ProductUserResponse

	query := "UPDATE products SET category_id = ?, product_name = ?, color = ?, stock = ?, price = ? WHERE id = ?"

	if err := db.DB.Exec(query, product.CategoryID, product.ProductName, product.Color, product.Stock, product.Price, product.ID).Error; err != nil {
		return models.ProductUserResponse{}, err
	}

	if err := db.DB.Raw("select * from products where id = ?", product.ID).Scan(&modProduct).Error; err != nil {
		return models.ProductUserResponse{}, err
	}
	return modProduct, nil
}
func (i *ProductRepository) DeleteProduct(productID string) error {

	id, err := strconv.Atoi(productID)
	if err != nil {
		return errors.New("converting into integet is not happened")
	}

	result := i.DB.Exec("DELETE FROM products WHERE id = ?", id)

	if result.RowsAffected < 1 {
		return errors.New("no records with that ID exist")
	}

	return nil
}

func (i *ProductRepository) DeleteProducts(inventoryID string) error {
	id, err := strconv.Atoi(inventoryID)
	if err != nil {
		return errors.New("converting into interger is not happened")
	}
	result := i.DB.Exec("DELETE FROM products WHERE id=?", id)
	// fmt.Println("ID fr repo", id)
	// fmt.Println("result error", result.Error)
	// fmt.Println("rows affected", result.RowsAffected)
	if result.RowsAffected < 1 {
		return errors.New("no records with that ID exist")

	}
	return nil
}
func (i *ProductRepository) CheckProducts(pid int) (bool, error) {
	var k int
	err := i.DB.Raw("SELECT COUNT(*) FROM products WHERE id=?", pid).Scan(&k).Error
	if err != nil {
		return false, err
	}
	if k == 0 {
		return false, err
	}
	return true, err
}

func (i *ProductRepository) UpdateProducts(pid int, stock int) (models.ProductResponse, error) {
	//Check the Database Connection
	if i.DB == nil {
		return models.ProductResponse{}, errors.New("database connection is nil")
	}

	//Update the stock
	if err := i.DB.Exec("UPDATE products SET stock=$1 WHERE id= $2", stock, pid).Error; err != nil {
		return models.ProductResponse{}, err

	}

	//Retrive the update
	var newDetails models.ProductResponse
	if err := i.DB.Raw("SELECT * FROM products WHERE id =? ", pid).Scan(&newDetails).Error; err != nil {
		return models.ProductResponse{}, err
	}

	return newDetails, nil

}
func (cr *ProductRepository) CheckProductAvailable(product_id int) (bool, error) {
	var count int
	querry := "SELECT COUNT(*) FROM products where id=?"
	err := cr.DB.Raw(querry, product_id).Scan(&count).Error
	if err != nil {
		return false, errors.New("product does not exist")
	}
	if count < 1 {
		return false, errors.New("product doesnot exist")
	}
	return true, nil
}

func (cr *ProductRepository) GetPriceOfProduct(product_id int) (float64, error) {
	querry := "SELECT price FROM products where id=?"
	var price float64
	err := cr.DB.Raw(querry, product_id).Scan(&price).Error

	if err != nil {
		return 0, errors.New("error in getting price")
	}
	return price, nil

}
