package merchantsrv

import (
	"database/sql"
	"errors"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/vanneeza/go-mnc/internal/domain/entity"
	"github.com/vanneeza/go-mnc/internal/domain/web/customerweb"
	"github.com/vanneeza/go-mnc/internal/domain/web/merchantweb"
	"github.com/vanneeza/go-mnc/internal/domain/web/productweb"
	"github.com/vanneeza/go-mnc/internal/domain/web/txweb"
	"github.com/vanneeza/go-mnc/internal/repository/customerrepo"
	"github.com/vanneeza/go-mnc/internal/repository/logrepo"
	"github.com/vanneeza/go-mnc/internal/repository/merchantrepo"
	"github.com/vanneeza/go-mnc/internal/repository/productrepo"

	"github.com/vanneeza/go-mnc/utils/helper"
	"github.com/vanneeza/go-mnc/utils/log"
	"golang.org/x/crypto/bcrypt"
)

type MerchantServiceImpl struct {
	Db                 *sql.DB
	MerchantRepository merchantrepo.MerchantRepository
	ProductRepository  productrepo.ProductRepository
	CustomerRepository customerrepo.CustomerRepository
	LogRepository      logrepo.LogRepository
}

// MerchantTxHistory implements TxService.
func (service *MerchantServiceImpl) MerchantTxHistory(merchantId string) ([]txweb.TxHistoryMerchantResponse, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	merchantTxHistory, err2 := service.MerchantRepository.HistoryTransaction(tx, merchantId)
	helper.PanicError(err2)

	var responses []txweb.TxHistoryMerchantResponse

	log := entity.Log{
		UserId:      merchantId,
		Level:       "info",
		Activity:    "Get",
		Description: "View history transaction",
		CreatedAt:   time.Now(),
	}
	service.LogRepository.Create(tx, log)
	for _, v := range merchantTxHistory {

		product, err6 := service.ProductRepository.FindByParams(tx, v.Product.Id, "")
		helper.PanicError(err6)

		customer, err4 := service.CustomerRepository.FindByParams(tx, v.Customer.Id, "")
		helper.PanicError(err4)

		detailResponse := txweb.TxDetailMerchant{
			Id:         v.Detail.Id,
			Status:     v.Detail.Status,
			TotalPrice: v.Detail.TotalPrice,
			Photo:      v.Detail.Photo,
		}

		productResponse := productweb.ProductMerchant{
			Id:          product.Id,
			Name:        product.Name,
			Price:       product.Price,
			Description: product.Description,
		}

		customerResponse := customerweb.Response{
			Id:      customer.Id,
			Name:    customer.Name,
			Phone:   customer.Phone,
			Address: customer.Address,
		}

		Response := txweb.TxHistoryMerchantResponse{
			Detail:   detailResponse,
			Product:  productResponse,
			Customer: customerResponse,
		}

		responses = append(responses, Response)
	}

	return responses, nil
}

// LogActivity implements MerchantService.
func (service *MerchantServiceImpl) LogActivity(merchantId string) (*[]entity.LogResponse, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	logs, err2 := service.LogRepository.GetAllByUserId(tx, merchantId)
	helper.PanicError(err2)

	var logResponse []entity.LogResponse
	for _, l := range *logs {
		logRsp := entity.LogResponse{
			Activity:    l.Activity,
			Description: l.Description,
			CreatedAt:   time.Now(),
		}
		logResponse = append(logResponse, logRsp)
	}

	return &logResponse, nil
}

// Edit implements MerchantService.
func (service *MerchantServiceImpl) Edit(req *merchantweb.UpdateRequest) (*merchantweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	merchant := entity.Merchant{
		Id:       req.Id,
		Name:     req.Name,
		Phone:    req.Phone,
		Password: string(encryptedPassword),
	}
	c, err2 := service.MerchantRepository.Update(tx, &merchant)
	if err2 != nil {
		return &merchantweb.Response{}, errors.New("failed update data merchant")
	} else {
		log := entity.Log{
			UserId:      merchant.Id,
			Level:       "info",
			Activity:    "Edit",
			Description: merchant.Name + " Have edit profile",
			CreatedAt:   time.Now(),
		}
		service.LogRepository.Create(tx, log)
		merchantResponse := merchantweb.Response{
			Id:    c.Id,
			Name:  c.Name,
			Phone: c.Phone,
		}
		return &merchantResponse, nil
	}

}

// GetAll implements MerchantService.
func (service *MerchantServiceImpl) GetAll() ([]merchantweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)

	defer helper.CommitOrRollback(tx)

	var merchants []merchantweb.Response

	merchant, err2 := service.MerchantRepository.FindAll(tx)
	if err2 != nil {
		return []merchantweb.Response{}, errors.New("data merchant not found")
	} else {
		for _, c := range merchant {
			merchantResponse := merchantweb.Response{
				Id:    c.Id,
				Name:  c.Name,
				Phone: c.Phone,
			}
			merchants = append(merchants, merchantResponse)
		}

		return merchants, nil
	}

}

// GetByParams implements MerchantService.
func (service *MerchantServiceImpl) GetByParams(merchantId, phone string) (*merchantweb.Response, *merchantweb.ForLogin, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	merchant, err2 := service.MerchantRepository.FindByParams(tx, merchantId, phone)
	if err2 != nil {
		log.Error("userID: "+merchantId, err)
		return nil, nil, err2
	}

	log := entity.Log{
		UserId:      merchant.Id,
		Level:       "info",
		Activity:    "View Profile",
		Description: merchant.Name + " Has Viewed the Profile",
		CreatedAt:   time.Now(),
	}
	service.LogRepository.Create(tx, log)

	merchantResponse := merchantweb.Response{
		Id:    merchant.Id,
		Name:  merchant.Name,
		Phone: merchant.Phone,
	}

	LoginResponse := merchantweb.ForLogin{
		Id:       merchant.Id,
		Name:     merchant.Name,
		Phone:    merchant.Phone,
		Password: merchant.Password,
		Role:     merchant.Role,
	}

	return &merchantResponse, &LoginResponse, nil
}

// Register implements MerchantService.
func (service *MerchantServiceImpl) Register(req *merchantweb.CreateRequest) (*merchantweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	merchant := entity.Merchant{
		Id:       ksuid.New().String(),
		Name:     req.Name,
		Phone:    req.Phone,
		Password: string(encryptedPassword),
	}

	c, err2 := service.MerchantRepository.Create(tx, &merchant)
	helper.PanicError(err2)

	log := entity.Log{
		UserId:      c.Id,
		Level:       "Info",
		Activity:    "Register",
		Description: "Register Merchant With Name " + c.Name,
		CreatedAt:   time.Now(),
	}
	service.LogRepository.Create(tx, log)

	merchantResponse := merchantweb.Response{
		Id:    c.Id,
		Name:  c.Name,
		Phone: c.Phone,
	}

	return &merchantResponse, nil
}

// Unreg implements MerchantService.
func (service *MerchantServiceImpl) Unreg(merchantId string) (*merchantweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	c, err2 := service.MerchantRepository.FindByParams(tx, merchantId, "")
	helper.PanicError(err2)

	log := entity.Log{
		UserId:      c.Id,
		Level:       "Info",
		Activity:    "Delete",
		Description: c.Name + " Has Delete Account Merchant",
		CreatedAt:   time.Now(),
	}
	service.LogRepository.Create(tx, log)
	merchantResponse := merchantweb.Response{
		Id:    c.Id,
		Name:  c.Name,
		Phone: c.Phone,
	}

	service.MerchantRepository.Delete(tx, merchantId)
	return &merchantResponse, nil
}

func NewMerchantService(db *sql.DB, merchantRepository merchantrepo.MerchantRepository, customerRepository customerrepo.CustomerRepository, productRepository productrepo.ProductRepository, logRepository logrepo.LogRepository) MerchantService {
	return &MerchantServiceImpl{
		Db:                 db,
		MerchantRepository: merchantRepository,
		LogRepository:      logRepository,
		ProductRepository:  productRepository,
		CustomerRepository: customerRepository,
	}
}
