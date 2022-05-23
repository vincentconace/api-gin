package domain

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	//ID          *int     `json:"id"`
	ProductCode string  `json:"product_code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
}
