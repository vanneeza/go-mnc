package merchantsrv

import (
	"github.com/vanneeza/go-mnc/internal/domain/entity"
	"github.com/vanneeza/go-mnc/internal/domain/web/merchantweb"
	"github.com/vanneeza/go-mnc/internal/domain/web/txweb"
)

type MerchantService interface {
	Register(req *merchantweb.CreateRequest) (*merchantweb.Response, error)
	GetAll() ([]merchantweb.Response, error)
	GetByParams(merchantId, phone string) (*merchantweb.Response, *merchantweb.ForLogin, error)
	Edit(req *merchantweb.UpdateRequest) (*merchantweb.Response, error)
	Unreg(merchantId string) (*merchantweb.Response, error)
	LogActivity(merchantId string) (*[]entity.LogResponse, error)
	MerchantTxHistory(merchantId string) ([]txweb.TxHistoryMerchantResponse, error)
}
