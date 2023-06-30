package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vanneeza/go-mnc/api/controller/bankctrl"
	"github.com/vanneeza/go-mnc/api/controller/customerctrl"
	"github.com/vanneeza/go-mnc/api/controller/loginctrl"
	"github.com/vanneeza/go-mnc/api/controller/merchantctrl"
	"github.com/vanneeza/go-mnc/api/controller/productctrl"
	"github.com/vanneeza/go-mnc/api/controller/txctrl"
	"github.com/vanneeza/go-mnc/internal/middleware"
	"github.com/vanneeza/go-mnc/internal/repository/balancerepo"
	"github.com/vanneeza/go-mnc/internal/repository/bankrepo"
	"github.com/vanneeza/go-mnc/internal/repository/customerrepo"
	"github.com/vanneeza/go-mnc/internal/repository/logrepo"
	"github.com/vanneeza/go-mnc/internal/repository/merchantrepo"
	"github.com/vanneeza/go-mnc/internal/repository/productrepo"
	"github.com/vanneeza/go-mnc/internal/repository/txrepo"
	"github.com/vanneeza/go-mnc/internal/service/banksrv"
	"github.com/vanneeza/go-mnc/internal/service/customersrv"
	"github.com/vanneeza/go-mnc/internal/service/logsrv"
	"github.com/vanneeza/go-mnc/internal/service/merchantsrv"
	"github.com/vanneeza/go-mnc/internal/service/productsrv"
	"github.com/vanneeza/go-mnc/internal/service/txsrv"
)

func Run(db *sql.DB, jwtKey string) *gin.Engine {
	r := gin.Default()

	logRepo := logrepo.NewLogRepository()
	logSrv := logsrv.NewLogService(db, logRepo)
	validate := validator.New()

	productRepo := productrepo.NewProductRepository()
	customerRepo := customerrepo.NewCustomerRepository()
	merchantRepo := merchantrepo.NewMerchantRepository()
	bankRepo := bankrepo.NewBankRepository()
	balanceRepo := balancerepo.NewBalanceRepository()

	customerSrv := customersrv.NewCustomerService(db, customerRepo, merchantRepo, productRepo, logRepo)
	customerCtrl := customerctrl.NewCustomerController(customerSrv)

	merchantSrv := merchantsrv.NewMerchantService(db, merchantRepo, customerRepo, productRepo, logRepo)
	merchantCtrl := merchantctrl.NewMerchantController(merchantSrv)

	bankSrv := banksrv.NewBankService(db, bankRepo, merchantRepo, logSrv)
	bankCtrl := bankctrl.NewBankController(bankSrv)

	productSrv := productsrv.NewProductService(db, productRepo, merchantRepo, logSrv)
	productCtrl := productctrl.NewProductController(productSrv)

	txRepo := txrepo.NewTxRepository()
	txSrv := txsrv.NewTxService(db, txRepo, productRepo, merchantRepo, bankRepo, balanceRepo, customerRepo, validate, logSrv)
	txCtrl := txctrl.NewTxController(txSrv)

	loginCtrl := loginctrl.NewLoginController(customerSrv, merchantSrv, validate, logSrv)
	api := r.Group("mnc/api/")
	api.POST("login/", loginCtrl.Login)
	api.GET("transaction/orders", txCtrl.ViewAllOrder)
	api.POST("transaction/confirmation/:id", txCtrl.Confirmation)

	api.POST("customer/", customerCtrl.Register)
	api.GET("customers/", customerCtrl.GetAll)
	api.POST("merchant/", merchantCtrl.Register)
	api.GET("merchants/", merchantCtrl.GetAll)
	api.GET("merchant/banks", bankCtrl.GetAll)
	api.GET("products", productCtrl.GetAll)

	api.POST("admin/bank", bankCtrl.RegisterBankAdmin)
	api.GET("admin/banks", bankCtrl.GetAllBankAdmin)
	api.GET("admin/bank/:id", bankCtrl.GetByIdBankAdmin)

	customer := api.Group("customer")
	customer.Use(middleware.AuthMiddleware(jwtKey))
	{

		customer.GET("/", customerCtrl.GetByParams)
		customer.PUT("/", customerCtrl.Edit)
		customer.DELETE("/", customerCtrl.Unreg)

		customer.GET("activity", customerCtrl.LogActivity)
		customer.GET("transaction_history", customerCtrl.ViewTxHistories)
		customer.POST("logout/", loginCtrl.Logout)
	}

	merchant := api.Group("merchant")
	merchant.Use(middleware.AuthMiddleware(jwtKey))
	{
		merchant.GET("/", merchantCtrl.GetByParams)
		merchant.PUT("/", merchantCtrl.Edit)
		merchant.DELETE("/", merchantCtrl.Unreg)

		merchant.POST("bank/", bankCtrl.Register)
		merchant.GET("bank/:id", bankCtrl.GetByParams)
		merchant.PUT("bank/:id", bankCtrl.Edit)
		merchant.DELETE("bank/:id", bankCtrl.Unreg)

		merchant.POST("product/", productCtrl.Register)
		merchant.GET("product/:id", productCtrl.GetByParams)
		merchant.PUT("product/:id", productCtrl.Edit)
		merchant.DELETE("product/:id", productCtrl.Unreg)

		merchant.GET("activity", merchantCtrl.LogActivity)
		merchant.GET("transaction_history", merchantCtrl.HistoryTransaction)
		merchant.POST("logout/", loginCtrl.Logout)
	}

	tx := api.Group("transaction")

	tx.Use(middleware.AuthMiddleware(jwtKey))
	{
		tx.POST("/order", txCtrl.Invoice)
		tx.POST("/payment/:id", txCtrl.Payment)

	}

	return r
}
