package product

import (
	"context"

	"github.com/vincentconace/api-gin/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Get(ctx context.Context) ([]domain.Product, error)
	GetById(ctx context.Context, id uint) (domain.Product, error)
	Save(ctx context.Context, p domain.Product) (uint, error)
	Update(ctx context.Context, id uint, p domain.Product) (domain.Product, error)
	Delete(ctx context.Context, id uint) error
	Exists(ctx context.Context, productCode string) bool
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Query all products
// var (
// 	getProductsQuery    = `SELECT id, product_code, name, description, price, stock FROM products`
// 	getProductByIdQuery = `SELECT id, product_code, name, description, price, stock FROM products WHERE id = ?`
// 	createProductQuery  = `INSERT INTO products (product_code, name, description, price, stock) VALUES (?, ?, ?, ?, ?)`
// 	updateProductQuery  = `UPDATE products SET product_code = ?, name = ?, description = ?, price = ?, stock = ? WHERE id = ?`
// 	deleteProductQuery  = `DELETE FROM products WHERE id = ?`
// 	existProductQuery   = `SELECT id FROM products WHERE product_code = ?`
// )

func (r *repository) Get(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	result := r.db.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (r *repository) GetById(ctx context.Context, id uint) (domain.Product, error) {
	var p domain.Product
	result := r.db.First(&p, id).Error
	if result != nil {
		return domain.Product{}, result
	}
	return p, nil
}

func (r *repository) Save(ctx context.Context, p domain.Product) (uint, error) {
	result := r.db.Create(&p).Error
	if result != nil {
		return 0, result
	}
	return p.ID, nil
}

func (r *repository) Update(ctx context.Context, id uint, p domain.Product) (domain.Product, error) {
	var product domain.Product
	result := r.db.Model(&product).Where("id = ?", id).Updates(p).Error

	if result != nil {
		return domain.Product{}, result
	}

	product.ID = id

	return product, nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	result := r.db.Delete(&domain.Product{}, id).Error
	if result != nil {
		return result
	}

	return nil
}

func (r *repository) Exists(ctx context.Context, productCode string) bool {
	result := r.db.Where("product_code = ?", productCode).First(&domain.Product{}).Error
	return result == nil
}
