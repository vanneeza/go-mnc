package productsrv

import (
	"database/sql"

	"github.com/segmentio/ksuid"
	"github.com/vanneeza/go-mnc/internal/domain/entity"
	"github.com/vanneeza/go-mnc/internal/domain/web/merchantweb"
	"github.com/vanneeza/go-mnc/internal/domain/web/productweb"
	"github.com/vanneeza/go-mnc/internal/repository/merchantrepo"
	"github.com/vanneeza/go-mnc/internal/repository/productrepo"
	"github.com/vanneeza/go-mnc/internal/service/logsrv"

	"github.com/vanneeza/go-mnc/utils/helper"
)

type ProductServiceImpl struct {
	Db                 *sql.DB
	ProductRepository  productrepo.ProductRepository
	MerchantRepository merchantrepo.MerchantRepository
	logsrv             logsrv.LogService
}

// Edit implements ProductService.
func (service *ProductServiceImpl) Edit(req *productweb.UpdateRequest) (*productweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	merchant, err3 := service.MerchantRepository.FindByParams(tx, req.Merchant.Id, "")
	helper.PanicError(err3)

	product := entity.Product{
		Id:          req.Id,
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Merchant:    *merchant,
	}
	p, err2 := service.ProductRepository.Update(tx, &product)
	helper.PanicError(err2)

	log := entity.Log{
		UserId:      merchant.Id,
		Activity:    "Edit",
		Description: merchant.Name + " has edit data product with name " + product.Name,
	}
	service.logsrv.Register(&log)

	merchantResponse := merchantweb.Response{
		Id:    merchant.Id,
		Name:  merchant.Name,
		Phone: merchant.Phone,
	}

	productResponse := productweb.Response{
		Id:          p.Id,
		Name:        p.Name,
		Price:       p.Price,
		Description: p.Description,
		Merchant:    merchantResponse,
	}
	return &productResponse, nil
}

// GetAll implements ProductService.
func (service *ProductServiceImpl) GetAll() ([]productweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)

	defer helper.CommitOrRollback(tx)

	var products []productweb.Response
	product, err2 := service.ProductRepository.FindAll(tx)
	helper.PanicError(err2)

	for _, p := range product {
		merchant, err3 := service.MerchantRepository.FindByParams(tx, p.Merchant.Id, "")
		helper.PanicError(err3)

		merchantResponse := merchantweb.Response{
			Id:    merchant.Id,
			Name:  merchant.Name,
			Phone: merchant.Phone,
		}

		productResponse := productweb.Response{
			Id:          p.Id,
			Name:        p.Name,
			Price:       p.Price,
			Description: p.Description,
			Merchant:    merchantResponse,
		}
		products = append(products, productResponse)
	}

	return products, nil

}

// GetByParams implements ProductService.
func (service *ProductServiceImpl) GetByParams(productId, merchantId string) (*productweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	product, err2 := service.ProductRepository.FindByParams(tx, productId, merchantId)
	helper.PanicError(err2)

	merchant, err3 := service.MerchantRepository.FindByParams(tx, product.Merchant.Id, "")
	helper.PanicError(err3)

	log := entity.Log{
		UserId:      merchant.Id,
		Activity:    "Get",
		Description: merchant.Name + " has view the product with name " + product.Name,
	}
	service.logsrv.Register(&log)

	merchantResponse := merchantweb.Response{
		Id:    merchant.Id,
		Name:  merchant.Name,
		Phone: merchant.Phone,
	}

	productResponse := productweb.Response{
		Id:          product.Id,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Merchant:    merchantResponse,
	}

	return &productResponse, nil
}

// Register implements ProductService.
func (service *ProductServiceImpl) Register(req *productweb.CreateRequest) (*productweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	merchant, err3 := service.MerchantRepository.FindByParams(tx, req.Merchant.Id, "")
	helper.PanicError(err3)

	product := entity.Product{
		Id:          ksuid.New().String(),
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Merchant:    *merchant,
	}

	p, err2 := service.ProductRepository.Create(tx, &product)
	helper.PanicError(err2)

	log := entity.Log{
		UserId:      merchant.Id,
		Activity:    "Add",
		Description: merchant.Name + " has add the product with name " + product.Name,
	}
	service.logsrv.Register(&log)

	merchantResponse := merchantweb.Response{
		Id:    merchant.Id,
		Name:  merchant.Name,
		Phone: merchant.Phone,
	}

	productResponse := productweb.Response{
		Id:          p.Id,
		Name:        p.Name,
		Price:       p.Price,
		Description: p.Description,
		Merchant:    merchantResponse,
	}
	return &productResponse, nil
}

// Unreg implements ProductService.
func (service *ProductServiceImpl) Unreg(productId string) (*productweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	p, err2 := service.ProductRepository.FindByParams(tx, productId, "")
	helper.PanicError(err2)

	merchant, err3 := service.MerchantRepository.FindByParams(tx, p.Merchant.Id, "")
	helper.PanicError(err3)

	log := entity.Log{
		UserId:      merchant.Id,
		Activity:    "Delete",
		Description: merchant.Name + " has delete the product with name " + p.Name,
	}
	service.logsrv.Register(&log)

	merchantResponse := merchantweb.Response{
		Id:    merchant.Id,
		Name:  merchant.Name,
		Phone: merchant.Phone,
	}
	productResponse := productweb.Response{
		Id:          p.Id,
		Name:        p.Name,
		Price:       p.Price,
		Description: p.Description,
		Merchant:    merchantResponse,
	}

	service.ProductRepository.Delete(tx, productId)

	return &productResponse, nil
}

func NewProductService(db *sql.DB, productRepository productrepo.ProductRepository, merchantRepository merchantrepo.MerchantRepository, logService logsrv.LogService) ProductService {
	return &ProductServiceImpl{
		Db:                 db,
		ProductRepository:  productRepository,
		MerchantRepository: merchantRepository,
		logsrv:             logService,
	}
}
