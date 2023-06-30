package customerrepo

import (
	"database/sql"

	"github.com/vanneeza/go-mnc/internal/domain/entity"
)

type CustomerRepository interface {
	Create(tx *sql.Tx, customer *entity.Customer) (*entity.Customer, error)
	FindAll(tx *sql.Tx) ([]entity.Customer, error)
	FindByParams(tx *sql.Tx, customerId, phone string) (*entity.Customer, error)
	Update(tx *sql.Tx, customer *entity.Customer) (*entity.Customer, error)
	Delete(tx *sql.Tx, customerId string) error
	HistoryTransaction(tx *sql.Tx, customerId string) ([]entity.CustomerTxHistory, error)
}
