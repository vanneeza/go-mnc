package bankrepo

import (
	"database/sql"

	"github.com/vanneeza/go-mnc/internal/domain/entity"
)

type BankRepository interface {
	Create(tx *sql.Tx, bank *entity.Bank) (*entity.Bank, error)
	FindAll(tx *sql.Tx) ([]entity.Bank, error)
	FindByParams(tx *sql.Tx, bankId, merchant_id string) ([]entity.Bank, error)
	Update(tx *sql.Tx, bank *entity.Bank) (*entity.Bank, error)
	Delete(tx *sql.Tx, bankId string) error

	CreateBankAdmin(tx *sql.Tx, bank *entity.BankAdmin) (*entity.BankAdmin, error)
	FindAllBankAdmin(tx *sql.Tx) ([]entity.BankAdmin, error)
	FindByIdBankAdmin(tx *sql.Tx, bankId string) (*entity.BankAdmin, error)
}
