package banksrv

import "github.com/vanneeza/go-mnc/internal/domain/web/bankweb"

type BankService interface {
	Register(req *bankweb.CreateRequest) (*bankweb.Response, error)
	GetAll() ([]bankweb.Response, error)
	GetByParams(bankId, merchantId string) ([]bankweb.Response, error)
	Edit(req *bankweb.UpdateRequest) (*bankweb.Response, error)
	Unreg(bankId string) (*[]bankweb.Response, error)

	RegisterBankAdmin(req *bankweb.BankAdminCreateRequest) (*bankweb.ResponseForDetail, error)
	GetAllBankAdmin() ([]bankweb.ResponseForDetail, error)
	GetByIdBankAdmin(bankId string) (bankweb.ResponseForDetail, error)
}
