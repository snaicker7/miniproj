package repository

import (
	"database/sql"
	"miniproj/entity"
)

type DatabaseRepoInterface interface {
	Connection() *sql.DB
	CreateUserDetails() error
	UpdateProductTable(product *entity.Product) error
	InsertProductTable(users *entity.UserAPIs) (int, error)
	FetchProductImgUrl(productId int) (*entity.Product, error)
	FetchProductDetails(productId int) (*entity.Product, error)
}
