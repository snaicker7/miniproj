package entity

import "time"

type UserAPIs struct {
	UserId             int      `json:"user_id" required:"true"`
	ProductName        string   `json:"product_name"`
	ProductDescription string   `json:"product_description"`
	ProductImages      []string `json:"product_images"`
	ProductPrice       float32  `json:"product_price"`
}
type Product struct {
	ProductID               int       `json:"product_id" required:"true"`
	ProductName             string    `json:"product_name"`
	ProductDescription      string    `json:"product_description"`
	ProductImages           []string  `json:"product_images"`
	CompressedProductImages []string  `json:"compressed_product_images"`
	ProductPrice            float64   `json:"product_price"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}
