package balancerepo

import (
	"database/sql"
	"errors"

	"github.com/vanneeza/go-mnc/internal/domain/entity"
	"github.com/vanneeza/go-mnc/utils/helper"
)

type BalanceRepositoryImpl struct {
}

func NewBalanceRepository() BalanceRepository {
	return &BalanceRepositoryImpl{}
}

func (repository *BalanceRepositoryImpl) FindAll(tx *sql.Tx) (*entity.Balance, error) {
	SQL := "SELECT id, balance FROM balance"
	rows, err := tx.Query(SQL)
	helper.PanicError(err)
	defer rows.Close()

	var balance entity.Balance
	if rows.Next() {
		err := rows.Scan(&balance.Id, &balance.Balance)
		helper.PanicError(err)
		return &balance, nil
	} else {
		return nil, errors.New("balance not found")
	}
}

func (repository *BalanceRepositoryImpl) Update(tx *sql.Tx, balance *entity.Balance) (*entity.Balance, error) {
	SQL := "UPDATE balance SET balance = $1 WHERE id = $2"
	_, err := tx.Exec(SQL, balance.Balance, balance.Id)
	helper.PanicError(err)

	return balance, nil
}
