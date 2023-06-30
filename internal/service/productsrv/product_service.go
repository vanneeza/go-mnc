package productsrv

import "github.com/vanneeza/go-mnc/internal/domain/web/productweb"

type ProductService interface {
	Register(req *productweb.CreateRequest) (*productweb.Response, error)
	GetAll() ([]productweb.Response, error)
	GetByParams(productId, merchantId string) (*productweb.Response, error)
	Edit(req *productweb.UpdateRequest) (*productweb.Response, error)
	Unreg(productId string) (*productweb.Response, error)
}
