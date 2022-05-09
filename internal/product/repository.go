package product

import (
	"context"
	"database/sql"

	"github.com/vincentconace/api-gin/internal/domain"
)

type Repository interface {
	Get(ctx context.Context) ([]domain.Product, error)
	GetById(ctx context.Context, id int) (domain.Product, error)
	Save(ctx context.Context, p domain.Product) (int, error)
	Update(ctx context.Context, id int, p domain.Product) error
	Delete(ctx context.Context, id int) error
	Exists(ctx context.Context, productCode string) bool
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// Query all products
var (
	getProductsQuery    = `SELECT id, product_code, name, description, price, stock FROM products`
	getProductByIdQuery = `SELECT id, product_code, name, description, price, stock FROM products WHERE id = ?`
	createProductQuery  = `INSERT INTO products (product_code, name, description, price, stock) VALUES (?, ?, ?, ?, ?)`
	updateProductQuery  = `UPDATE products SET product_code = ?, name = ?, description = ?, price = ?, stock = ? WHERE id = ?`
	deleteProductQuery  = `DELETE FROM products WHERE id = ?`
	existProductQuery   = `SELECT id FROM products WHERE product_code = ?`
)

func (r *repository) Get(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	rows, err := r.db.Query(getProductsQuery)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.ProductCode, &p.Name, &p.Description, &p.Price, &p.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *repository) GetById(ctx context.Context, id int) (domain.Product, error) {
	var p domain.Product
	err := r.db.QueryRow(getProductByIdQuery, id).Scan(&p.ID, &p.ProductCode, &p.Name, &p.Description, &p.Price, &p.Stock)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (r *repository) Save(ctx context.Context, p domain.Product) (int, error) {
	stmt, err := r.db.Prepare(createProductQuery)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(p.ProductCode, p.Name, p.Description, p.Price, p.Stock)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *repository) Update(ctx context.Context, id int, p domain.Product) error {
	stmt, err := r.db.Prepare(updateProductQuery)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(&p.ProductCode, &p.Name, &p.Description, &p.Price, &p.Stock, id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	stmt, err := r.db.Prepare(deleteProductQuery)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows < 1 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) Exists(ctx context.Context, productCode string) bool {
	row := r.db.QueryRow(existProductQuery, productCode)
	err := row.Scan(&productCode)
	return err == nil
}
