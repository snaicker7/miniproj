package dbrepo

import (
	"database/sql"
	"miniproj/entity"
)

type TestDBRepo struct{}

func (m *TestDBRepo) Connection() *sql.DB {
	return nil
}

func (m *TestDBRepo) CreateUserDetails() error {
	return nil
}

func (m *TestDBRepo) UpdateProductTable(product *entity.Product) error {

	return nil
}

func (m *TestDBRepo) InsertProductTable(users entity.UserAPIs) (int, error) {
	var product_id int = 1
	return product_id, nil
}

func (m *TestDBRepo) FetchProductImgUrl(productId int) (*entity.Product, error) {
	var product = entity.Product{
		ProductID:          1,
		ProductName:        "ProductName",
		ProductDescription: "ProductDescription",
		ProductImages:      []string{"img1.jpg", "img2.jpg"},
	}
	return &product, nil
}
func (m *TestDBRepo) FetchProductDetails(productId int) (*entity.Product, error) {
	var product = entity.Product{
		ProductID:          1,
		ProductName:        "ProductName",
		ProductDescription: "ProductDescription",
		ProductImages:      []string{"img1.jpg", "img2.jpg"},
	}
	return &product, nil
}
