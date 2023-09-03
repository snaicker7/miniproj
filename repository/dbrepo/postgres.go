package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"miniproj/entity"
	"time"

	"github.com/lib/pq"
)

const (
	dbTimeout = time.Second * 3
)

type PostgresDatabaseRepo struct {
	DB *sql.DB
}

func (m *PostgresDatabaseRepo) Connection() *sql.DB {
	return m.DB
}

func (p *PostgresDatabaseRepo) CreateUserDetails() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	usersInsertQuery := `INSERT INTO "users" (id, name, mobile, longitude, latitude, created_at, updated_at) VALUES (1,'Surya Naicker', '1234567890',-74.0060,40.7128,'2023-08-26 10:00:00', '2023-08-26 10:00:00'), (2,'Iron Man', '9876543210',-118.2437,34.0522, '2023-08-26 10:15:00', '2023-08-26 10:15:00'), (3,'Barbie', '5555566666',-0.1278,51.5074, '2023-08-26 10:30:00', '2023-08-26 10:30:00');`
	_, err := p.DB.ExecContext(ctx, usersInsertQuery)
	if err != nil {
		log.Printf("Error %s when inserting values in users table\n", err)
		return err
	}
	return nil
}

func (p *PostgresDatabaseRepo) UpdateProductTable(product *entity.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	updateSmt := `UPDATE products SET compressed_product_images = $1, updated_at = NOW() WHERE product_id = $2;`
	_, err := p.DB.ExecContext(ctx, updateSmt, pq.Array(product.CompressedProductImages), product.ProductID)
	if err != nil {
		log.Printf("Error %s when updating compressed_product_images:%v for product_id:% v in products table\n", err, product.ProductImages, product.ProductID)
		return err
	}
	return nil
}

func (p *PostgresDatabaseRepo) InsertProductTable(users *entity.UserAPIs) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var product_id int
	insertDynStmt := `INSERT INTO products (product_name, product_description, product_images, product_price, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING product_id`
	err := p.DB.QueryRowContext(ctx, insertDynStmt,
		users.ProductName,
		users.ProductDescription,
		pq.Array(users.ProductImages),
		users.ProductPrice).Scan(&product_id)
	if err != nil {
		log.Printf("Error %s when inserting details in products table\n", err)
		return 0, err
	}
	fmt.Println(product_id)
	return product_id, nil
}

func (p *PostgresDatabaseRepo) FetchProductImgUrl(productId int) (*entity.Product, error) {
	query := `select product_id,product_name,product_description, product_images, product_price, created_at from products where product_id=$1`
	var product entity.Product
	err := p.DB.QueryRow(query, productId).Scan(
		&product.ProductID,
		&product.ProductName,
		&product.ProductDescription,
		pq.Array(&product.ProductImages),
		&product.ProductPrice,
		&product.CreatedAt)
	if err != nil {
		errSmt := fmt.Sprintf("Error %s when fetching product_name and product_images user table\n", err)
		log.Fatal(errSmt)
		return nil, errors.New(errSmt)
	}
	return &product, nil
}

func (p *PostgresDatabaseRepo) FetchProductDetails(productId int) (*entity.Product, error) {
	query := `select product_id,product_name,product_description, product_images,compressed_product_images, product_price, created_at, updated_at from products where product_id=$1`
	var product entity.Product
	err := p.DB.QueryRow(query, productId).Scan(
		&product.ProductID,
		&product.ProductName,
		&product.ProductDescription,
		pq.Array(&product.ProductImages),
		pq.Array(&product.CompressedProductImages),
		&product.ProductPrice,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		errSmt := fmt.Sprintf("Error %s when fetching product_name and product_images user table\n", err)
		log.Fatal(errSmt)
		return nil, errors.New(errSmt)
	}
	return &product, nil
}
