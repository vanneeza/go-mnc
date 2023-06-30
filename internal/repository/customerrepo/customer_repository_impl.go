package customerrepo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/vanneeza/go-mnc/internal/domain/entity"
	"github.com/vanneeza/go-mnc/utils/helper"
)

type CustomerRepositoryImpl struct {
}

func NewCustomerRepository() CustomerRepository {
	return &CustomerRepositoryImpl{}
}

func (repository *CustomerRepositoryImpl) Create(tx *sql.Tx, customer *entity.Customer) (*entity.Customer, error) {

	SQL := "INSERT INTO customer(id, name, phone, address, password) VALUES($1, $2, $3, $4, $5)"
	_, err := tx.Exec(SQL, customer.Id, customer.Name, customer.Phone, customer.Address, customer.Password)
	helper.PanicError(err)
	return customer, nil
}

func (repository *CustomerRepositoryImpl) FindAll(tx *sql.Tx) ([]entity.Customer, error) {
	SQL := "SELECT id, name, phone, address, password FROM customer WHERE is_deleted = false"
	rows, err := tx.Query(SQL)
	helper.PanicError(err)

	var categories []entity.Customer
	for rows.Next() {
		customer := entity.Customer{}
		err := rows.Scan(&customer.Id, &customer.Name, &customer.Phone, &customer.Address, &customer.Password)
		helper.PanicError(err)

		categories = append(categories, customer)
	}

	defer rows.Close()
	return categories, nil
}
func (repository *CustomerRepositoryImpl) FindByParams(tx *sql.Tx, customerId, phone string) (*entity.Customer, error) {
	SQL := "SELECT id, name, phone, address, password, role FROM customer WHERE is_deleted = false AND id = $1 OR phone = $2;"
	rows, err := tx.Query(SQL, customerId, phone)
	helper.PanicError(err)

	defer rows.Close()

	customer := entity.Customer{}
	if rows.Next() {
		err := rows.Scan(&customer.Id, &customer.Name, &customer.Phone, &customer.Address, &customer.Password, &customer.Role)
		helper.PanicError(err)
		fmt.Printf("repo customer: %v\n", customer)
		return &customer, nil
	} else {
		return &customer, errors.New("customer is not found")
	}
}

func (repository *CustomerRepositoryImpl) Update(tx *sql.Tx, customer *entity.Customer) (*entity.Customer, error) {
	SQL := "UPDATE customer SET name = $1, phone = $2, address = $3, password = $4 WHERE id = $5"
	_, err := tx.Exec(SQL, customer.Name, customer.Phone, customer.Address, customer.Password, customer.Id)
	helper.PanicError(err)

	return customer, nil
}

func (repository *CustomerRepositoryImpl) Delete(tx *sql.Tx, customerId string) error {
	SQL := "UPDATE customer SET is_deleted = TRUE WHERE id = $1"
	_, err := tx.Exec(SQL, customerId)
	helper.PanicError(err)

	return nil
}

func (*CustomerRepositoryImpl) HistoryTransaction(tx *sql.Tx, customerId string) ([]entity.CustomerTxHistory, error) {
	SQL := `SELECT detail.id, detail.status, detail.total_price, detail.photo,
	tx_order.qty, tx_order.product_id, product.merchant_id FROM detail
	INNER JOIN tx_order ON detail.id = tx_order.detail_id
	INNER JOIN product ON tx_order.product_id = product.id
	WHERE tx_order.customer_id = $1`
	rows, err := tx.Query(SQL, customerId)
	helper.PanicError(err)
	defer rows.Close()

	var customerTxHistories []entity.CustomerTxHistory
	for rows.Next() {
		var customerTx entity.CustomerTxHistory
		err := rows.Scan(&customerTx.Detail.Id, &customerTx.Detail.Status, &customerTx.Detail.TotalPrice, &customerTx.Detail.Photo, &customerTx.Qty, &customerTx.Product.Id, &customerTx.Merchant.Id)
		helper.PanicError(err)

		customerTxHistories = append(customerTxHistories, customerTx)
	}

	return customerTxHistories, nil
}
