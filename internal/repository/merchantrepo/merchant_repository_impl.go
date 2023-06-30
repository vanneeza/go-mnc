package merchantrepo

import (
	"database/sql"
	"errors"

	"github.com/vanneeza/go-mnc/internal/domain/entity"
	"github.com/vanneeza/go-mnc/utils/helper"
)

type MerchantRepositoryImpl struct {
}

// HistoryTransaction implements MerchantRepository.
func (*MerchantRepositoryImpl) HistoryTransaction(tx *sql.Tx, merchantId string) ([]entity.MerchantTxHistory, error) {
	SQL := `SELECT detail.id, detail.status, detail.total_price, detail.photo,
	tx_order.qty, tx_order.product_id, tx_order.customer_id FROM detail
	INNER JOIN tx_order ON detail.id = tx_order.detail_id
	INNER JOIN product ON tx_order.product_id = product.id
	WHERE detail.status = 'Paid' AND product.merchant_id = $1`
	rows, err := tx.Query(SQL, merchantId)
	helper.PanicError(err)
	defer rows.Close()

	var MerchantTxHistories []entity.MerchantTxHistory
	for rows.Next() {
		var merchantTx entity.MerchantTxHistory
		err := rows.Scan(&merchantTx.Detail.Id, &merchantTx.Detail.Status, &merchantTx.Detail.TotalPrice, &merchantTx.Detail.Photo, &merchantTx.Qty, &merchantTx.Product.Id, &merchantTx.Customer.Id)
		helper.PanicError(err)

		MerchantTxHistories = append(MerchantTxHistories, merchantTx)
	}

	return MerchantTxHistories, nil
}

func NewMerchantRepository() MerchantRepository {
	return &MerchantRepositoryImpl{}
}

func (repository *MerchantRepositoryImpl) Create(tx *sql.Tx, merchant *entity.Merchant) (*entity.Merchant, error) {

	SQL := "INSERT INTO merchant(id, name, phone, password) VALUES($1, $2, $3, $4)"
	_, err := tx.Exec(SQL, merchant.Id, merchant.Name, merchant.Phone, merchant.Password)
	helper.PanicError(err)
	return merchant, nil
}

func (repository *MerchantRepositoryImpl) FindAll(tx *sql.Tx) ([]entity.Merchant, error) {
	SQL := "SELECT id, name, phone, password FROM merchant WHERE is_deleted = false"
	rows, err := tx.Query(SQL)
	helper.PanicError(err)

	var categories []entity.Merchant
	for rows.Next() {
		merchant := entity.Merchant{}
		err := rows.Scan(&merchant.Id, &merchant.Name, &merchant.Phone, &merchant.Password)
		helper.PanicError(err)

		categories = append(categories, merchant)
	}

	defer rows.Close()
	return categories, nil
}
func (repository *MerchantRepositoryImpl) FindByParams(tx *sql.Tx, merchantId, phone string) (*entity.Merchant, error) {
	SQL := "SELECT id, name, phone, password, role FROM merchant WHERE is_deleted = false AND id = $1 OR phone = $2;"
	rows, err := tx.Query(SQL, merchantId, phone)
	helper.PanicError(err)

	defer rows.Close()

	merchant := entity.Merchant{}
	if rows.Next() {
		err := rows.Scan(&merchant.Id, &merchant.Name, &merchant.Phone, &merchant.Password, &merchant.Role)
		helper.PanicError(err)

		return &merchant, nil
	} else {
		return &merchant, errors.New("merchant is not found")
	}
}

func (repository *MerchantRepositoryImpl) Update(tx *sql.Tx, merchant *entity.Merchant) (*entity.Merchant, error) {
	SQL := "UPDATE merchant SET name = $1, phone = $2, password = $3 WHERE id = $4"
	_, err := tx.Exec(SQL, merchant.Name, merchant.Phone, merchant.Password, merchant.Id)
	helper.PanicError(err)

	return merchant, nil
}

func (repository *MerchantRepositoryImpl) Delete(tx *sql.Tx, merchantId string) error {
	SQL := "UPDATE merchant SET is_deleted = TRUE WHERE id = $1"
	_, err := tx.Exec(SQL, merchantId)
	helper.PanicError(err)

	return nil
}
