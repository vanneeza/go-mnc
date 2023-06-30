package txrepo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/vanneeza/go-mnc/internal/domain/entity"
	"github.com/vanneeza/go-mnc/utils/helper"
)

type TxRepositoryImpl struct {
}

// GetAllOrder implements TxRepository.
func (*TxRepositoryImpl) GetAllOrder(tx *sql.Tx, detailId, status string) ([]entity.OrderDetail, error) {
	SQL := `
	SELECT tx_order.id, tx_order.qty, tx_order.product_id, tx_order.customer_id, detail.id AS detail_id,
	detail.status, detail.total_price, detail.photo, detail.bank_id, payment.pay
	FROM detail
	INNER JOIN tx_order ON detail.id = tx_order.detail_id
	INNER JOIN payment ON detail.id = payment.detail_id
	WHERE detail.id = $1 OR detail.status = $2	
	`
	rows, err := tx.Query(SQL, detailId, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.OrderDetail
	for rows.Next() {
		var od entity.OrderDetail
		err := rows.Scan(
			&od.Order.Id,
			&od.Order.Qty,
			&od.Order.Product.Id,
			&od.Order.Customer.Id,
			&od.Detail.Id,
			&od.Detail.Status,
			&od.Detail.TotalPrice,
			&od.Detail.Photo,
			&od.Detail.Bank.Id,
			&od.Pay,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, od)
	}

	return orders, nil
}

// FindDetail implements TxRepository.
func (*TxRepositoryImpl) FindDetail(tx *sql.Tx, detailId string) (*entity.Detail, error) {
	SQL := "SELECT id, status, total_price, photo, bank_id FROM detail WHERE detail_id = $1"
	rows, err := tx.Query(SQL, detailId)
	helper.PanicError(err)
	defer rows.Close()

	var detail entity.Detail
	if rows.Next() {
		err := rows.Scan(&detail.Id, &detail.Status, &detail.TotalPrice, &detail.Photo, &detail.Bank.Id)
		helper.PanicError(err)
		return &detail, nil

	} else {
		return &detail, errors.New("detail is not found")
	}
}

// ConfirmationOrder implements TxRepository.
func (*TxRepositoryImpl) ConfirmationOrder(tx *sql.Tx, detail *entity.Detail) (*entity.Detail, error) {
	SQL := "UPDATE detail SET status = $1 WHERE id = $2"
	_, err := tx.Exec(SQL, detail.Status, detail.Id)
	helper.PanicError(err)

	return detail, nil
}

// CreateDetail implements TxRepository.
func (repository *TxRepositoryImpl) CreateDetail(tx *sql.Tx, detail *entity.Detail) (*entity.Detail, error) {
	SQL := "INSERT INTO detail(id, status, total_price, bank_id, photo) VALUES($1, $2, $3, $4, $5)"
	_, err := tx.Exec(SQL, detail.Id, detail.Status, detail.TotalPrice, detail.Bank.Id, detail.Photo)
	helper.PanicError(err)

	return detail, nil
}

// CreateOrder implements TxRepository.
func (repository *TxRepositoryImpl) CreateOrder(tx *sql.Tx, order *entity.Order) (*entity.Order, error) {
	fmt.Printf("repo order: %v\n", order)
	SQL := "INSERT INTO tx_order(id, qty, product_id, customer_id, detail_id) VALUES($1, $2, $3, $4, $5)"
	_, err := tx.Exec(SQL, order.Id, order.Qty, order.Product.Id, order.Customer.Id, order.Detail.Id)
	helper.PanicError(err)

	return order, nil
}

// CreatePayment implements TxRepository.
func (repository *TxRepositoryImpl) CreatePayment(tx *sql.Tx, payment *entity.Payment) (*entity.Payment, error) {
	SQL := "INSERT INTO payment(id, pay, detail_id) VALUES($1, $2, $3)"
	_, err := tx.Exec(SQL, payment.Id, payment.Pay, payment.Detail.Id)
	helper.PanicError(err)

	return payment, nil
}

// UpdateStatus implements TxRepository.
func (repository *TxRepositoryImpl) UpdateDetail(tx *sql.Tx, detail *entity.Detail) (*entity.Detail, error) {
	SQL := "UPDATE detail SET status = $1, photo = $2, bank_id = $3 WHERE id = $4"
	_, err := tx.Exec(SQL, detail.Status, detail.Photo, detail.Bank.Id, detail.Id)
	helper.PanicError(err)

	return detail, nil
}

func (repository *TxRepositoryImpl) FindOrder(tx *sql.Tx, detailId string) (*entity.Order, error) {
	SQL := `SELECT tx_order.id, tx_order.qty, tx_order.product_id, tx_order.customer_id, tx_order.detail_id,
	detail.status, detail.total_price, detail.bank_id, detail.photo FROM detail
	INNER JOIN tx_order ON detail.id = tx_order.detail_id WHERE detail.id = $1`
	rows, err := tx.Query(SQL, detailId)
	helper.PanicError(err)
	defer rows.Close()

	var order entity.Order
	if rows.Next() {
		err := rows.Scan(&order.Id, &order.Qty, &order.Product.Id, &order.Customer.Id, &order.Detail.Id, &order.Detail.Status, &order.Detail.TotalPrice, &order.Detail.Bank.Id, &order.Detail.Photo)
		helper.PanicError(err)
		return &order, nil

	} else {
		return &order, errors.New("order is not found")
	}

}

func NewTxRepository() TxRepository {
	return &TxRepositoryImpl{}
}
