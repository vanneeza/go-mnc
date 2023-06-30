package customersrv

import (
	"github.com/vanneeza/go-mnc/internal/domain/entity"
	"github.com/vanneeza/go-mnc/internal/domain/web/customerweb"
	"github.com/vanneeza/go-mnc/internal/domain/web/txweb"
)

type CustomerService interface {
	Register(req *customerweb.CreateRequest) (*customerweb.Response, error)
	GetAll() ([]customerweb.Response, error)
	GetByParams(customerId, phone string) (*customerweb.Response, *customerweb.ForLogin, error)
	Edit(req *customerweb.UpdateRequest) (*customerweb.Response, error)
	Unreg(customerId string) (*customerweb.Response, error)
	LogActivity(customerId string) (*[]entity.Log, error)
	ViewTxHistories(customerId string) ([]txweb.TxHistoryCustomerResponse, error)
}
