package customersrv

import (
	"database/sql"
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
	"golang.org/x/crypto/bcrypt"
)

type CustomerServiceImpl struct {
	Db                 *sql.DB
	CustomerRepository customerrepo.CustomerRepository
	MerchantRepository merchantrepo.MerchantRepository
	ProductRepository  productrepo.ProductRepository
	LogRepository      logrepo.LogRepository
}

// ViewTxHistories implements CustomerService.
func (service *CustomerServiceImpl) ViewTxHistories(customerId string) ([]txweb.TxHistoryCustomerResponse, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	customerTxHistory, err2 := service.CustomerRepository.HistoryTransaction(tx, customerId)
	helper.PanicError(err2)

	var responses []txweb.TxHistoryCustomerResponse

	log := entity.Log{
		UserId:      customerId,
		Level:       "info",
		Activity:    "Get",
		Description: "View history transaction",
		CreatedAt:   time.Now(),
	}
	service.LogRepository.Create(tx, log)

	for _, v := range customerTxHistory {

		product, err6 := service.ProductRepository.FindByParams(tx, v.Product.Id, "")
		helper.PanicError(err6)

		merchant, err4 := service.MerchantRepository.FindByParams(tx, v.Merchant.Id, "")
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

		merchantResponse := merchantweb.Response{
			Id:    merchant.Id,
			Name:  merchant.Name,
			Phone: merchant.Phone,
		}

		Response := txweb.TxHistoryCustomerResponse{
			Detail:   detailResponse,
			Product:  productResponse,
			Merchant: merchantResponse,
		}

		responses = append(responses, Response)
	}

	return responses, nil
}

// LogActivity implements CustomerService.
func (service *CustomerServiceImpl) LogActivity(customerId string) (*[]entity.Log, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	logs, err2 := service.LogRepository.GetAllByUserId(tx, customerId)
	helper.PanicError(err2)

	var logResponse []entity.Log

	for _, l := range *logs {
		logResponse = append(logResponse, l)
	}

	return &logResponse, nil
}

// Edit implements CustomerService.
func (service *CustomerServiceImpl) Edit(req *customerweb.UpdateRequest) (*customerweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	customer := entity.Customer{
		Id:       req.Id,
		Name:     req.Name,
		Phone:    req.Phone,
		Address:  req.Address,
		Password: string(encryptedPassword),
	}
	c, err2 := service.CustomerRepository.Update(tx, &customer)
	helper.PanicError(err2)

	log := entity.Log{
		UserId:      customer.Id,
		Level:       "Info",
		Activity:    "Edit",
		Description: customer.Name + " has edit data profile",
		CreatedAt:   time.Now(),
	}
	service.LogRepository.Create(tx, log)

	customerResponse := customerweb.Response{
		Id:      c.Id,
		Name:    c.Name,
		Phone:   c.Phone,
		Address: c.Address,
	}
	return &customerResponse, nil
}

// GetAll implements CustomerService.
func (service *CustomerServiceImpl) GetAll() ([]customerweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)

	defer helper.CommitOrRollback(tx)

	var customers []customerweb.Response
	customer, err2 := service.CustomerRepository.FindAll(tx)
	helper.PanicError(err2)

	for _, c := range customer {
		customerResponse := customerweb.Response{
			Id:      c.Id,
			Name:    c.Name,
			Phone:   c.Phone,
			Address: c.Address,
		}
		customers = append(customers, customerResponse)
	}

	return customers, nil

}

// GetByParams implements CustomerService.
func (service *CustomerServiceImpl) GetByParams(customerId, phone string) (*customerweb.Response, *customerweb.ForLogin, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	customer, err2 := service.CustomerRepository.FindByParams(tx, customerId, phone)
	helper.PanicError(err2)

	log := entity.Log{
		UserId:      customer.Id,
		Level:       "Info",
		Activity:    "Get",
		Description: customer.Name + " has view the profile",
		CreatedAt:   time.Now(),
	}
	service.LogRepository.Create(tx, log)

	customerResponse := customerweb.Response{
		Id:      customer.Id,
		Name:    customer.Name,
		Phone:   customer.Phone,
		Address: customer.Address,
	}

	LoginResponse := customerweb.ForLogin{
		Id:       customer.Id,
		Name:     customer.Name,
		Phone:    customer.Phone,
		Address:  customer.Address,
		Password: customer.Password,
		Role:     customer.Role,
	}

	return &customerResponse, &LoginResponse, nil
}

// Register implements CustomerService.
func (service *CustomerServiceImpl) Register(req *customerweb.CreateRequest) (*customerweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	customer := entity.Customer{
		Id:       ksuid.New().String(),
		Name:     req.Name,
		Phone:    req.Phone,
		Address:  req.Address,
		Password: string(encryptedPassword),
	}

	c, err2 := service.CustomerRepository.Create(tx, &customer)
	helper.PanicError(err2)

	log := entity.Log{
		UserId:      customer.Id,
		Level:       "Info",
		Activity:    "Create",
		Description: customer.Name + " has create new account customer ",
		CreatedAt:   time.Now(),
	}
	service.LogRepository.Create(tx, log)

	customerResponse := customerweb.Response{
		Id:      c.Id,
		Name:    c.Name,
		Phone:   c.Phone,
		Address: c.Address,
	}

	return &customerResponse, nil
}

// Unreg implements CustomerService.
func (service *CustomerServiceImpl) Unreg(customerId string) (*customerweb.Response, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	c, err2 := service.CustomerRepository.FindByParams(tx, customerId, "")
	helper.PanicError(err2)
	customerResponse := customerweb.Response{
		Id:      c.Id,
		Name:    c.Name,
		Phone:   c.Phone,
		Address: c.Address,
	}

	service.CustomerRepository.Delete(tx, customerId)
	log := entity.Log{
		UserId:      c.Id,
		Level:       "Info",
		Activity:    "Delete",
		Description: c.Name + " has delete account",
		CreatedAt:   time.Now(),
	}
	service.LogRepository.Create(tx, log)
	return &customerResponse, nil
}

func NewCustomerService(db *sql.DB, customerRepository customerrepo.CustomerRepository, merchantRepository merchantrepo.MerchantRepository, productRepository productrepo.ProductRepository, logRepo logrepo.LogRepository) CustomerService {
	return &CustomerServiceImpl{
		Db:                 db,
		CustomerRepository: customerRepository,
		MerchantRepository: merchantRepository,
		ProductRepository:  productRepository,
		LogRepository:      logRepo,
	}
}
