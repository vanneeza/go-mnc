package balancerepo

import (
	"database/sql"

	"github.com/vanneeza/go-mnc/internal/domain/entity"
)

type BalanceRepository interface {
	Update(tx *sql.Tx, balance *entity.Balance) (*entity.Balance, error)
	FindAll(tx *sql.Tx) (*entity.Balance, error)
}
