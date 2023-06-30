package merchantrepo

import (
	"database/sql"

	"github.com/vanneeza/go-mnc/internal/domain/entity"
)

type MerchantRepository interface {
	Create(tx *sql.Tx, merchant *entity.Merchant) (*entity.Merchant, error)
	FindAll(tx *sql.Tx) ([]entity.Merchant, error)
	FindByParams(tx *sql.Tx, merchantId, phone string) (*entity.Merchant, error)
	Update(tx *sql.Tx, merchant *entity.Merchant) (*entity.Merchant, error)
	Delete(tx *sql.Tx, merchantId string) error
	HistoryTransaction(tx *sql.Tx, merchantId string) ([]entity.MerchantTxHistory, error)
}
