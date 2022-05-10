package product

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vincentconace/api-gin/internal/domain"
)

func puntInt(i int) *int {
	return &i
}

func puntFloat(i float32) *float32 {
	return &i
}

func puntStr(i string) *string {
	return &i
}

var ctx = context.Background()

var productMock = []domain.Product{
	{
		ID:          puntInt(1),
		Name:        puntStr("Product 1"),
		Description: puntStr("Product 1 description"),
		Price:       puntFloat(1.99),
		Stock:       puntInt(10),
	},
	{
		ID:    puntInt(2),
		Name:  puntStr("Product 2"),
		Price: puntFloat(2.99),
		Stock: puntInt(20),
	},
}

func TestSaveOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO products")
	mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

	repository := NewRepository(db)
	idResult, err := repository.Save(ctx, productMock[0])

	assert.NoError(t, err)
	assert.Equal(t, *productMock[0].ID, idResult)
}

func TestSaveErrCreated(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO products")
	mock.ExpectExec("").WillReturnError(ErrCreatedProduct)

	repository := NewRepository(db)
	idResult, err := repository.Save(ctx, productMock[0])

	assert.EqualError(t, err, ErrCreatedProduct.Error())
	assert.Empty(t, idResult)
}

func TestGetOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	defer db.Close()

	colums := []string{"id", "product_code", "name", "description", "price", "stock"}

	rows := mock.NewRows(colums)
	rows.AddRow(1, "PRO001", "one", "description", 20, 50).AddRow(2, "PRO001", "two", "description", 20, 50)

	mock.ExpectQuery("SELECT id, product_code, name, description, price, stock FROM products").WillReturnRows(rows)

	repository := NewRepository(db)
	products, err := repository.Get(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
}

func TestGetErrNotConnectedDataBase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	defer db.Close()

	colums := []string{"id", "product_code", "name", "description", "price", "stock"}
	rows := mock.NewRows(colums)
	rows.AddRow(1, "PRO001", "Product 1", "Product 1 description", 1.99, 10).AddRow(2, "PRO002", "Product 2", "Product 2 description", 2.99, 20)

	mock.ExpectQuery("SELECT id, product_code, name, description, price, stock FROM products").WillReturnError(sql.ErrConnDone)

	repository := NewRepository(db)
	products, err := repository.Get(ctx)

	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Empty(t, products)
}

func TestGetByIdOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	defer db.Close()

	colums := []string{"id", "product_code", "name", "description", "price", "stock"}
	rows := mock.NewRows(colums)
	rows.AddRow(1, "PRO001", "Product 1", "Product 1 description", 1.99, 10)

	mock.ExpectQuery("SELECT id, product_code, name, description, price, stock FROM products WHERE id = ?").WillReturnRows(rows).WithArgs(1)

	repository := NewRepository(db)
	product, err := repository.GetById(ctx, *productMock[0].ID)

	assert.NoError(t, err)
	assert.Equal(t, *productMock[0].ID, *product.ID)
	assert.Equal(t, *productMock[0].Name, *product.Name)
	assert.Equal(t, *productMock[0].Description, *product.Description)
	assert.Equal(t, *productMock[0].Price, *product.Price)
	assert.Equal(t, *productMock[0].Stock, *product.Stock)
}

func TestGetByIdErrNotFount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	defer db.Close()

	colums := []string{"id", "product_code", "name", "description", "price", "stock"}
	rows := mock.NewRows(colums)
	rows.AddRow(1, "PRO001", "Product 1", "Product 1 description", 1.99, 10)

	mock.ExpectQuery("SELECT id, product_code, name, description, price, stock FROM products WHERE id = ?").WillReturnError(ErrNotFound)

	repository := NewRepository(db)
	product, err := repository.GetById(ctx, *productMock[0].ID)

	assert.EqualError(t, err, ErrNotFound.Error())
	assert.Empty(t, product)
}

func TestUpdateOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	defer db.Close()

	mock.ExpectPrepare("UPDATE products")
	mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

	repository := NewRepository(db)
	err = repository.Update(ctx, 1, productMock[0])

	assert.NoError(t, err)
}

func TestUpdateErrNotFount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	defer db.Close()

	mock.ExpectPrepare("UPDATE products")
	mock.ExpectExec("").WillReturnError(ErrNotFound)

	repository := NewRepository(db)
	err = repository.Update(ctx, 1, productMock[0])

	assert.EqualError(t, err, ErrNotFound.Error())
}

/*func TestExistOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	defer db.Close()

	colums := []string{"id", "product_code", "name", "description", "price", "stock"}
	rows := mock.NewRows(colums)
	rows.AddRow(1, "PRO001", "Product 1", "Product 1 description", 1.99, 10)

	mock.ExpectQuery("SELECT id FROM products WHERE product_code = ?").WithArgs("PRO001").WillReturnRows(rows)

	repository := NewRepository(db)
	result := repository.Exists(ctx, "PRO001")

	assert.True(t, result)
}*/
