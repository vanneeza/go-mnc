package productrepo

import (
	"database/sql"
	"errors"

	"github.com/vanneeza/go-mnc/internal/domain/entity"
	"github.com/vanneeza/go-mnc/utils/helper"
)

type ProductRepositoryImpl struct {
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{}
}

func (repository *ProductRepositoryImpl) Create(tx *sql.Tx, product *entity.Product) (*entity.Product, error) {

	SQL := "INSERT INTO product(id, name, price, description, merchant_id) VALUES($1, $2, $3, $4, $5)"
	_, err := tx.Exec(SQL, product.Id, product.Name, product.Price, product.Description, product.Merchant.Id)
	helper.PanicError(err)
	return product, nil
}

func (repository *ProductRepositoryImpl) FindAll(tx *sql.Tx) ([]entity.Product, error) {
	SQL := "SELECT id, name, price, description, merchant_id FROM product WHERE is_deleted = false"
	rows, err := tx.Query(SQL)
	helper.PanicError(err)

	var products []entity.Product
	for rows.Next() {
		product := entity.Product{}
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Description, &product.Merchant.Id)
		helper.PanicError(err)

		products = append(products, product)
	}

	defer rows.Close()
	return products, nil
}
func (repository *ProductRepositoryImpl) FindByParams(tx *sql.Tx, productId, merchantId string) (*entity.Product, error) {
	SQL := "SELECT id, name, price, description, merchant_id FROM product WHERE is_deleted = false AND id = $1 OR merchant_id = $2;"
	rows, err := tx.Query(SQL, productId, merchantId)
	helper.PanicError(err)

	defer rows.Close()

	product := entity.Product{}
	if rows.Next() {
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Description, &product.Merchant.Id)
		helper.PanicError(err)
		return &product, nil
	} else {
		return &product, errors.New("product is not found")
	}
}

func (repository *ProductRepositoryImpl) Update(tx *sql.Tx, product *entity.Product) (*entity.Product, error) {
	SQL := "UPDATE product SET name = $1, price = $2, description = $3, merchant_id = $4 WHERE id = $5"
	_, err := tx.Exec(SQL, product.Name, product.Price, product.Description, product.Merchant.Id, product.Id)
	helper.PanicError(err)

	return product, nil
}

func (repository *ProductRepositoryImpl) Delete(tx *sql.Tx, productId string) error {
	SQL := "UPDATE product SET is_deleted = TRUE WHERE id = $1"
	_, err := tx.Exec(SQL, productId)
	helper.PanicError(err)

	return nil
}
