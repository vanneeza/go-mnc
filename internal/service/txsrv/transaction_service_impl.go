package txsrv

import (
	"database/sql"

	"github.com/segmentio/ksuid"
	"github.com/vanneeza/go-mnc/internal/domain/entity"
	"github.com/vanneeza/go-mnc/internal/domain/web/bankweb"
	"github.com/vanneeza/go-mnc/internal/domain/web/customerweb"
	"github.com/vanneeza/go-mnc/internal/domain/web/merchantweb"
	"github.com/vanneeza/go-mnc/internal/domain/web/productweb"
	"github.com/vanneeza/go-mnc/internal/domain/web/txweb"
	"github.com/vanneeza/go-mnc/internal/repository/balancerepo"
	"github.com/vanneeza/go-mnc/internal/repository/bankrepo"
	"github.com/vanneeza/go-mnc/internal/repository/customerrepo"
	"github.com/vanneeza/go-mnc/internal/repository/merchantrepo"
	"github.com/vanneeza/go-mnc/internal/repository/productrepo"
	"github.com/vanneeza/go-mnc/internal/repository/txrepo"
	"github.com/vanneeza/go-mnc/internal/service/logsrv"
	"github.com/vanneeza/go-mnc/utils/helper"
)

type TxServiceImpl struct {
	Db                 *sql.DB
	TxRepository       txrepo.TxRepository
	ProductRepository  productrepo.ProductRepository
	MerchantRepository merchantrepo.MerchantRepository
	BankRepository     bankrepo.BankRepository
	BalanceRepository  balancerepo.BalanceRepository
	CustomerRepository customerrepo.CustomerRepository
	LogService         logsrv.LogService
}

// ViewAllPayment implements TxService.
func (service *TxServiceImpl) ViewAllPayment() ([]txweb.OrderDetail, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	orderDetails, err2 := service.TxRepository.GetAllOrder(tx, "", "Done")
	helper.PanicError(err2)

	var responses []txweb.OrderDetail
	for _, v := range orderDetails {
		product, err6 := service.ProductRepository.FindByParams(tx, v.Order.Product.Id, "")
		helper.PanicError(err6)

		merchant, err7 := service.MerchantRepository.FindByParams(tx, product.Merchant.Id, "")
		helper.PanicError(err7)

		banks, err8 := service.BankRepository.FindByIdBankAdmin(tx, v.Detail.Bank.Id)
		helper.PanicError(err8)

		customer, err4 := service.CustomerRepository.FindByParams(tx, v.Order.Customer.Id, "")
		helper.PanicError(err4)

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

		customerResponse := customerweb.Response{
			Id:      customer.Id,
			Name:    customer.Name,
			Phone:   customer.Phone,
			Address: customer.Address,
		}

		bankResponse := bankweb.ResponseForDetail{
			Id:            banks.Id,
			Name:          banks.Name,
			BankAccount:   banks.BankAccount,
			Branch:        banks.Branch,
			AccountNumber: banks.AccountNumber,
		}

		orderResponse := txweb.OrderResponseWithoutDetail{
			Id:       v.Order.Id,
			Qty:      v.Order.Qty,
			Product:  productResponse,
			Customer: customerResponse,
		}

		response := txweb.OrderDetail{
			Id:         v.Detail.Id,
			Status:     v.Detail.Status,
			TotalPrice: v.Detail.TotalPrice,
			Pay:        v.Pay,
			Bank:       bankResponse,
			Photo:      v.Detail.Photo,
			Order:      orderResponse,
		}
		responses = append(responses, response)

	}

	return responses, nil
}

// Confirmation implements TxService.
func (service *TxServiceImpl) Confirmation(req *txweb.DetailUpdateRequest) ([]txweb.Confirmation, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	detail := entity.Detail{
		Id:     req.Id,
		Status: "Paid",
	}
	service.TxRepository.ConfirmationOrder(tx, &detail)

	orderDetails, err2 := service.TxRepository.GetAllOrder(tx, req.Id, "")
	helper.PanicError(err2)

	var responses []txweb.Confirmation

	for _, v := range orderDetails {
		product, err6 := service.ProductRepository.FindByParams(tx, v.Order.Product.Id, "")
		helper.PanicError(err6)

		merchant, err7 := service.MerchantRepository.FindByParams(tx, product.Merchant.Id, "")
		helper.PanicError(err7)

		banks, err8 := service.BankRepository.FindByIdBankAdmin(tx, v.Detail.Bank.Id)
		helper.PanicError(err8)

		bankMerchants, err9 := service.BankRepository.FindByParams(tx, "", merchant.Id)
		helper.PanicError(err9)

		var bankMerchant bankweb.ResponseForDetail
		for _, b := range bankMerchants {
			bankMerchant = bankweb.ResponseForDetail{
				Id:            b.Id,
				Name:          b.Name,
				BankAccount:   b.BankAccount,
				Branch:        b.Branch,
				AccountNumber: b.AccountNumber,
			}
		}

		customer, err4 := service.CustomerRepository.FindByParams(tx, v.Order.Customer.Id, "")
		helper.PanicError(err4)

		b, err3 := service.BalanceRepository.FindAll(tx)
		helper.PanicError(err3)

		balanceUpdate := b.Balance - v.Detail.TotalPrice
		balance := entity.Balance{
			Id:      b.Id,
			Balance: balanceUpdate,
		}

		service.BalanceRepository.Update(tx, &balance)

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

		customerResponse := customerweb.Response{
			Id:      customer.Id,
			Name:    customer.Name,
			Phone:   customer.Phone,
			Address: customer.Address,
		}

		bankResponse := bankweb.ResponseForDetail{
			Id:            banks.Id,
			Name:          banks.Name,
			BankAccount:   banks.BankAccount,
			Branch:        banks.Branch,
			AccountNumber: banks.AccountNumber,
		}

		orderResponse := txweb.OrderResponseWithoutDetail{
			Id:       v.Order.Id,
			Qty:      v.Order.Qty,
			Product:  productResponse,
			Customer: customerResponse,
		}

		payout := txweb.Payout{
			Payout:       v.Detail.TotalPrice,
			BankMerchant: bankMerchant,
		}

		response := txweb.Confirmation{
			Id:         v.Detail.Id,
			Status:     v.Detail.Status,
			TotalPrice: v.Detail.TotalPrice,
			Pay:        v.Pay,
			Bank:       bankResponse,
			Photo:      v.Detail.Photo,
			Order:      orderResponse,
			Payout:     payout,
		}

		responses = append(responses, response)

	}

	return responses, nil

}

// Invoice implements TxService.
func (service *TxServiceImpl) Invoice(req *txweb.OrderCreateRequest) (*txweb.OrderResponse, error) {

	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	customer, err3 := service.CustomerRepository.FindByParams(tx, req.CustomerId, "")
	helper.PanicError(err3)

	product, err4 := service.ProductRepository.FindByParams(tx, req.ProductId, "")
	helper.PanicError(err4)

	merchant, err5 := service.MerchantRepository.FindByParams(tx, product.Merchant.Id, "")
	helper.PanicError(err5)

	banks, err := service.BankRepository.FindAllBankAdmin(tx)
	helper.PanicError(err)

	banksResponse := make([]bankweb.ResponseForDetail, len(banks))
	var bank entity.BankAdmin
	for i, b := range banks {
		banksResponse[i] = bankweb.ResponseForDetail{
			Id:            b.Id,
			Name:          b.Name,
			BankAccount:   b.BankAccount,
			Branch:        b.Branch,
			AccountNumber: b.AccountNumber,
		}

		bank = entity.BankAdmin{
			Id:            b.Id,
			Name:          b.Name,
			BankAccount:   b.BankAccount,
			Branch:        b.Branch,
			AccountNumber: b.AccountNumber,
		}
	}

	totalPrice := product.Price * float64(req.Qty)

	detail := entity.Detail{
		Id:         ksuid.New().String(),
		Status:     "Pending",
		TotalPrice: totalPrice,
		Bank:       bank,
		Photo:      "NULL",
	}

	d, err2 := service.TxRepository.CreateDetail(tx, &detail)
	helper.PanicError(err2)

	detailResponse := txweb.DetailResponse{
		Id:         d.Id,
		Status:     d.Status,
		TotalPrice: totalPrice,
		Bank:       banksResponse,
		Photo:      d.Photo,
	}

	order := entity.Order{
		Id:       ksuid.New().String(),
		Qty:      req.Qty,
		Product:  *product,
		Customer: *customer,
		Detail:   *d,
	}
	o, err6 := service.TxRepository.CreateOrder(tx, &order)
	helper.PanicError(err6)

	log := entity.Log{
		UserId:      customer.Id,
		Activity:    "Create",
		Description: customer.Name + " has create the order with name product " + product.Name,
	}
	service.LogService.Register(&log)

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

	customerResponse := customerweb.Response{
		Id:      customer.Id,
		Name:    customer.Name,
		Phone:   customer.Phone,
		Address: customer.Address,
	}

	OrderResponse := txweb.OrderResponse{
		Id:       o.Id,
		Qty:      o.Qty,
		Product:  productResponse,
		Customer: customerResponse,
		Detail:   detailResponse,
	}
	return &OrderResponse, nil
}

// Payment implements TxService.
func (service *TxServiceImpl) Payment(req *txweb.PaymentCreateRequest) (*txweb.PaymentResponse, error) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	order, err5 := service.TxRepository.FindOrder(tx, req.DetailId)
	helper.PanicError(err5)

	product, err6 := service.ProductRepository.FindByParams(tx, order.Product.Id, "")
	helper.PanicError(err6)

	merchant, err7 := service.MerchantRepository.FindByParams(tx, product.Merchant.Id, "")
	helper.PanicError(err7)

	banks, err8 := service.BankRepository.FindByIdBankAdmin(tx, req.BankId)
	helper.PanicError(err8)

	customer, err4 := service.CustomerRepository.FindByParams(tx, req.CustomerId, "")
	helper.PanicError(err4)

	bankResponse := bankweb.ResponseForDetail{
		Id:            banks.Id,
		Name:          banks.Name,
		BankAccount:   banks.BankAccount,
		Branch:        banks.Branch,
		AccountNumber: banks.AccountNumber,
	}

	detail := entity.Detail{
		Id:     req.DetailId,
		Status: "Done",
		Bank:   *banks,
		Photo:  req.Photo.Filename,
	}
	payment := entity.Payment{
		Id:     ksuid.New().String(),
		Pay:    req.Pay,
		Detail: detail,
	}

	p, err2 := service.TxRepository.CreatePayment(tx, &payment)
	helper.PanicError(err2)

	_, err3 := service.TxRepository.UpdateDetail(tx, &payment.Detail)
	helper.PanicError(err3)

	balance, err9 := service.BalanceRepository.FindAll(tx)
	helper.PanicError(err9)

	totalPrice := product.Price * float64(order.Qty)
	sumBalance := balance.Balance + totalPrice

	updateBalance := entity.Balance{
		Id:      balance.Id,
		Balance: sumBalance,
	}

	service.BalanceRepository.Update(tx, &updateBalance)

	log := entity.Log{
		UserId:      customer.Id,
		Activity:    "Create",
		Description: customer.Name + " Has made payment on the order with ID order: " + order.Id,
	}
	service.LogService.Register(&log)

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

	customerResponse := customerweb.Response{
		Id:      customer.Id,
		Name:    customer.Name,
		Phone:   customer.Phone,
		Address: customer.Address,
	}

	orderResponse := txweb.PaymentOrder{
		Product:  productResponse,
		Customer: customerResponse,
		Bank:     bankResponse,
	}

	paymentResponse := txweb.PaymentResponse{
		Id:           p.Id,
		PaymentOrder: orderResponse,
		Pay:          req.Pay,
		Photo:        req.Photo.Filename,
	}
	return &paymentResponse, nil
}

func NewTxService(Db *sql.DB, txRepository txrepo.TxRepository, productRepository productrepo.ProductRepository, merchantRepository merchantrepo.MerchantRepository, bankRepository bankrepo.BankRepository, balanceRepository balancerepo.BalanceRepository, customerRepository customerrepo.CustomerRepository, logService logsrv.LogService) TxService {
	return &TxServiceImpl{
		Db:                 Db,
		TxRepository:       txRepository,
		ProductRepository:  productRepository,
		MerchantRepository: merchantRepository,
		BankRepository:     bankRepository,
		BalanceRepository:  balanceRepository,
		CustomerRepository: customerRepository,
		LogService:         logService,
	}
}
