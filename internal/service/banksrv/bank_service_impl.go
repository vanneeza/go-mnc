package banksrv

import (
	"database/sql"
	"fmt"

	"github.com/segmentio/ksuid"
	"github.com/vanneeza/go-mnc/internal/domain/entity"
	"github.com/vanneeza/go-mnc/internal/domain/web/bankweb"
	"github.com/vanneeza/go-mnc/internal/domain/web/merchantweb"
	"github.com/vanneeza/go-mnc/internal/repository/bankrepo"
	"github.com/vanneeza/go-mnc/internal/repository/merchantrepo"
	"github.com/vanneeza/go-mnc/internal/service/logsrv"

	"github.com/vanneeza/go-mnc/utils/helper"
)

type BankServiceImpl struct {
	Db                 *sql.DB
	BankRepository     bankrepo.BankRepository
	MerchantRepository merchantrepo.MerchantRepository
	LogService         logsrv.LogService
}

// GetByIdBankAdmin implements BankService.
func (service *BankServiceImpl) GetByIdBankAdmin(bankId string) (bankweb.ResponseForDetail, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	bank, err2 := service.BankRepository.FindByIdBankAdmin(tx, bankId)
	helper.PanicError(err2)

	bankResponse := bankweb.ResponseForDetail{
		Id:            bankId,
		Name:          bank.Name,
		BankAccount:   bank.BankAccount,
		Branch:        bank.Branch,
		AccountNumber: bank.AccountNumber,
	}

	return bankResponse, nil
}

// GetAllBankAdmin implements BankService.
func (service *BankServiceImpl) GetAllBankAdmin() ([]bankweb.ResponseForDetail, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)

	defer helper.CommitOrRollback(tx)

	var banks []bankweb.ResponseForDetail
	bank, err2 := service.BankRepository.FindAllBankAdmin(tx)
	helper.PanicError(err2)

	for _, b := range bank {
		bankResponse := bankweb.ResponseForDetail{
			Id:            b.Id,
			Name:          b.Name,
			BankAccount:   b.BankAccount,
			Branch:        b.Branch,
			AccountNumber: b.AccountNumber,
		}
		banks = append(banks, bankResponse)
	}

	return banks, nil
}

// RegisterBankAdmin implements BankService.
func (service *BankServiceImpl) RegisterBankAdmin(req *bankweb.BankAdminCreateRequest) (*bankweb.ResponseForDetail, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	bank := entity.BankAdmin{
		Id:            ksuid.New().String(),
		Name:          req.Name,
		BankAccount:   req.BankAccount,
		Branch:        req.Branch,
		AccountNumber: req.AccountNumber,
	}

	fmt.Printf("srv bank: %v\n", bank)

	b, err2 := service.BankRepository.CreateBankAdmin(tx, &bank)
	helper.PanicError(err2)

	bankResponse := bankweb.ResponseForDetail{
		Id:            b.Id,
		Name:          b.Name,
		BankAccount:   b.BankAccount,
		Branch:        b.Branch,
		AccountNumber: b.AccountNumber,
	}
	return &bankResponse, nil
}

// Edit implements BankService.
func (service *BankServiceImpl) Edit(req *bankweb.UpdateRequest) (*bankweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	merchant, err3 := service.MerchantRepository.FindByParams(tx, req.Merchant.Id, "")
	helper.PanicError(err3)

	bank := entity.Bank{
		Id:            req.Id,
		Name:          req.Name,
		BankAccount:   req.BankAccount,
		Branch:        req.Branch,
		AccountNumber: req.AccountNumber,
		Merchant:      *merchant,
	}

	b, err2 := service.BankRepository.Update(tx, &bank)
	helper.PanicError(err2)

	log := entity.Log{
		UserId:      merchant.Id,
		Activity:    "Edit",
		Description: merchant.Name + " Has edit data bank",
	}
	service.LogService.Register(&log)

	merchantResponse := merchantweb.Response{
		Id:    merchant.Id,
		Name:  merchant.Name,
		Phone: merchant.Phone,
	}

	bankResponse := bankweb.Response{
		Id:            b.Id,
		Name:          b.Name,
		BankAccount:   b.BankAccount,
		Branch:        b.Branch,
		AccountNumber: b.AccountNumber,
		Merchant:      merchantResponse,
	}
	return &bankResponse, nil
}

// GetAll implements BankService.
func (service *BankServiceImpl) GetAll() ([]bankweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)

	defer helper.CommitOrRollback(tx)

	var banks []bankweb.Response
	bank, err2 := service.BankRepository.FindAll(tx)
	helper.PanicError(err2)

	for _, b := range bank {
		merchant, err3 := service.MerchantRepository.FindByParams(tx, b.Merchant.Id, "")
		helper.PanicError(err3)

		merchantResponse := merchantweb.Response{
			Id:    merchant.Id,
			Name:  merchant.Name,
			Phone: merchant.Phone,
		}

		bankResponse := bankweb.Response{
			Id:            b.Id,
			Name:          b.Name,
			BankAccount:   b.BankAccount,
			Branch:        b.Branch,
			AccountNumber: b.AccountNumber,
			Merchant:      merchantResponse,
		}
		banks = append(banks, bankResponse)
	}

	return banks, nil

}

// GetByParams implements BankService.
func (service *BankServiceImpl) GetByParams(bankId, merchantId string) ([]bankweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	bank, err2 := service.BankRepository.FindByParams(tx, bankId, merchantId)
	helper.PanicError(err2)

	merchant, err3 := service.MerchantRepository.FindByParams(tx, merchantId, "")
	helper.PanicError(err3)

	log := entity.Log{
		UserId:      merchant.Id,
		Activity:    "Get",
		Description: merchant.Name + " Has get data bank by ID",
	}
	service.LogService.Register(&log)

	merchantResponse := merchantweb.Response{
		Id:    merchant.Id,
		Name:  merchant.Name,
		Phone: merchant.Phone,
	}

	banksResponse := make([]bankweb.Response, len(bank))
	for _, b := range bank {
		Response := bankweb.Response{
			Id:            b.Id,
			BankAccount:   b.BankAccount,
			Branch:        b.Branch,
			AccountNumber: b.AccountNumber,
			Merchant:      merchantResponse,
		}
		banksResponse = append(banksResponse, Response)
	}

	return banksResponse, nil
}

// Register implements BankService.
func (service *BankServiceImpl) Register(req *bankweb.CreateRequest) (*bankweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	merchant, err3 := service.MerchantRepository.FindByParams(tx, req.Merchant.Id, "")
	helper.PanicError(err3)

	bank := entity.Bank{
		Id:            ksuid.New().String(),
		Name:          req.Name,
		BankAccount:   req.BankAccount,
		Branch:        req.Branch,
		AccountNumber: req.AccountNumber,
		Merchant:      *merchant,
	}

	b, err2 := service.BankRepository.Create(tx, &bank)
	helper.PanicError(err2)

	log := entity.Log{
		UserId:      merchant.Id,
		Activity:    "Register",
		Description: merchant.Name + " Has add new data bank",
	}
	service.LogService.Register(&log)

	merchantResponse := merchantweb.Response{
		Id:    merchant.Id,
		Name:  merchant.Name,
		Phone: merchant.Phone,
	}

	bankResponse := bankweb.Response{
		Id:            b.Id,
		Name:          b.Name,
		BankAccount:   b.BankAccount,
		Branch:        b.Branch,
		AccountNumber: b.AccountNumber,
		Merchant:      merchantResponse,
	}
	return &bankResponse, nil
}

// Unreg implements BankService.
func (service *BankServiceImpl) Unreg(bankId string) (*[]bankweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	bank, err2 := service.BankRepository.FindByParams(tx, bankId, "")
	helper.PanicError(err2)

	var merchantId string
	for _, b := range bank {
		merchantId = b.Merchant.Id
	}

	merchant, err3 := service.MerchantRepository.FindByParams(tx, merchantId, "")
	helper.PanicError(err3)

	merchantResponse := merchantweb.Response{
		Id:    merchant.Id,
		Name:  merchant.Name,
		Phone: merchant.Phone,
	}
	banksResponse := make([]bankweb.Response, len(bank))

	for _, b := range bank {
		Response := bankweb.Response{
			Id:            b.Id,
			BankAccount:   b.BankAccount,
			Branch:        b.Branch,
			AccountNumber: b.AccountNumber,
			Merchant:      merchantResponse,
		}
		banksResponse = append(banksResponse, Response)
	}
	service.BankRepository.Delete(tx, bankId)

	log := entity.Log{
		UserId:      merchant.Id,
		Activity:    "Delete",
		Description: merchant.Name + " Has delete data bank",
	}
	service.LogService.Register(&log)

	return &banksResponse, nil
}

func NewBankService(db *sql.DB, bankRepository bankrepo.BankRepository, merchantRepository merchantrepo.MerchantRepository, logsrv logsrv.LogService) BankService {
	return &BankServiceImpl{
		Db:                 db,
		BankRepository:     bankRepository,
		MerchantRepository: merchantRepository,
		LogService:         logsrv,
	}
}
