package balancesrv

import (
	"database/sql"

	"github.com/vanneeza/go-mnc/internal/domain/entity"
	"github.com/vanneeza/go-mnc/internal/repository/balancerepo"
	"github.com/vanneeza/go-mnc/utils/helper"
)

type BalanceService interface {
	ViewAll() (*entity.BalanceResponse, error)
}

type BalanceServiceImpl struct {
	Db                *sql.DB
	BalanceRepository balancerepo.BalanceRepository
}

// ViewAll implements BalanceService.
func (service *BalanceServiceImpl) ViewAll() (*entity.BalanceResponse, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	b, _ := service.BalanceRepository.FindAll(tx)

	balanceResponse := entity.BalanceResponse{
		Balance: b.Balance,
	}
	return &balanceResponse, nil
}

func NewBalanceService(Db *sql.DB, balanceRepository balancerepo.BalanceRepository) BalanceService {
	return &BalanceServiceImpl{
		Db:                Db,
		BalanceRepository: balanceRepository,
	}
}
