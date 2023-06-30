package bankrepo

import (
	"database/sql"
	"errors"

	"github.com/vanneeza/go-mnc/internal/domain/entity"
	"github.com/vanneeza/go-mnc/utils/helper"
)

type BankRepositoryImpl struct {
}

func NewBankRepository() BankRepository {
	return &BankRepositoryImpl{}
}

func (*BankRepositoryImpl) FindByIdBankAdmin(tx *sql.Tx, bankId string) (*entity.BankAdmin, error) {
	SQL := "SELECT id, name, bank_account, branch, account_number FROM bank_admin WHERE id = $1;"
	rows, err := tx.Query(SQL, bankId)
	helper.PanicError(err)

	defer rows.Close()

	var bank entity.BankAdmin
	if rows.Next() {
		err := rows.Scan(&bank.Id, &bank.Name, &bank.BankAccount, &bank.Branch, &bank.AccountNumber)
		helper.PanicError(err)

		return &bank, nil
	} else {
		return &bank, errors.New("bank admin is not found")
	}
}

func (repository *BankRepositoryImpl) Create(tx *sql.Tx, bank *entity.Bank) (*entity.Bank, error) {

	SQL := "INSERT INTO bank(id, name, bank_account, branch, account_number, merchant_id) VALUES($1, $2, $3, $4, $5, $6)"
	_, err := tx.Exec(SQL, bank.Id, bank.Name, bank.BankAccount, bank.Branch, bank.AccountNumber, bank.Merchant.Id)
	helper.PanicError(err)
	return bank, nil
}

func (repository *BankRepositoryImpl) FindAll(tx *sql.Tx) ([]entity.Bank, error) {
	SQL := "SELECT id, name, bank_account, branch, account_number, merchant_id FROM bank"
	rows, err := tx.Query(SQL)
	helper.PanicError(err)

	var banks []entity.Bank
	for rows.Next() {
		bank := entity.Bank{}
		err := rows.Scan(&bank.Id, &bank.Name, &bank.BankAccount, &bank.Branch, &bank.AccountNumber, &bank.Merchant.Id)
		helper.PanicError(err)

		banks = append(banks, bank)
	}

	defer rows.Close()
	return banks, nil
}

func (repository *BankRepositoryImpl) FindByParams(tx *sql.Tx, bankId, merchant_id string) ([]entity.Bank, error) {
	SQL := "SELECT id, name, bank_account, branch, account_number, merchant_id FROM bank WHERE id = $1 OR merchant_id = $2;"
	rows, err := tx.Query(SQL, bankId, merchant_id)
	helper.PanicError(err)

	defer rows.Close()

	var banks []entity.Bank
	for rows.Next() {
		bank := entity.Bank{}
		err := rows.Scan(&bank.Id, &bank.Name, &bank.BankAccount, &bank.Branch, &bank.AccountNumber, &bank.Merchant.Id)
		helper.PanicError(err)

		banks = append(banks, bank)
	}

	if len(banks) > 0 {
		return banks, nil
	} else {
		return nil, errors.New("bank is not found")
	}
}

func (repository *BankRepositoryImpl) Update(tx *sql.Tx, bank *entity.Bank) (*entity.Bank, error) {
	SQL := "UPDATE bank SET name = $1, bank_account = $2, branch = $3, account_number = $4, merchant_id = $5 WHERE id = $6"
	_, err := tx.Exec(SQL, bank.Name, bank.BankAccount, bank.Branch, bank.AccountNumber, bank.Merchant.Id, bank.Id)
	helper.PanicError(err)

	return bank, nil
}

func (repository *BankRepositoryImpl) Delete(tx *sql.Tx, bankId string) error {
	SQL := "DELETE FROM bank WHERE id = $1"
	_, err := tx.Exec(SQL, bankId)
	helper.PanicError(err)

	return nil
}

// CreateBankAdmin implements BankRepository.
func (repository *BankRepositoryImpl) CreateBankAdmin(tx *sql.Tx, bank *entity.BankAdmin) (*entity.BankAdmin, error) {
	SQL := "INSERT INTO bank_admin(id, name, bank_account, branch, account_number) VALUES($1, $2, $3, $4, $5)"
	_, err := tx.Exec(SQL, bank.Id, bank.Name, bank.BankAccount, bank.Branch, bank.AccountNumber)
	helper.PanicError(err)
	return bank, nil
}

// FindAllBankAdmin implements BankRepository.
func (repository *BankRepositoryImpl) FindAllBankAdmin(tx *sql.Tx) ([]entity.BankAdmin, error) {
	SQL := "SELECT id, name, bank_account, branch, account_number FROM bank_admin"
	rows, err := tx.Query(SQL)
	helper.PanicError(err)

	var banks []entity.BankAdmin
	for rows.Next() {
		bank := entity.BankAdmin{}
		err := rows.Scan(&bank.Id, &bank.Name, &bank.BankAccount, &bank.Branch, &bank.AccountNumber)
		helper.PanicError(err)

		banks = append(banks, bank)
	}

	defer rows.Close()
	return banks, nil
}
