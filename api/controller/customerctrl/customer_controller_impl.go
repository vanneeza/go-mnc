package customerctrl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/vanneeza/go-mnc/internal/domain/web"
	"github.com/vanneeza/go-mnc/internal/domain/web/customerweb"
	"github.com/vanneeza/go-mnc/internal/service/customersrv"
	"github.com/vanneeza/go-mnc/utils/helper"
	"github.com/vanneeza/go-mnc/utils/log"
)

type CustomerControllerImpl struct {
	CustomerService customersrv.CustomerService
}

// ViewTxHistories implements CustomerController.
func (controller *CustomerControllerImpl) ViewTxHistories(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(jwt.MapClaims)
	customerId := claims["id"].(string)
	role := claims["role"].(string)

	if role != "customer" {
		result := web.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "unauthorized",
			Data:    "user is unauthorized",
		}
		ctx.JSON(http.StatusUnauthorized, result)
		return
	}

	response, err := controller.CustomerService.ViewTxHistories(customerId)
	if err != nil {
		log.Error("UserID: "+customerId, err)
		webResponse := web.WebResponse{
			Code:    http.StatusNotFound,
			Status:  "OK",
			Message: err.Error(),
			Data:    response,
		}

		ctx.JSON(http.StatusNotFound, gin.H{"transaction_history": webResponse})
	}

	log.Info("View transaction history customer id: " + customerId)
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "List of all transaction history",
		Data:    response,
	}

	ctx.JSON(http.StatusOK, gin.H{"transaction_history": webResponse})
}

// LogActivity implements CustomerController.
func (controller *CustomerControllerImpl) LogActivity(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(jwt.MapClaims)
	customerId := claims["id"].(string)
	role := claims["role"].(string)

	if role != "customer" {
		result := web.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "unauthorized",
			Data:    "user is unauthorized",
		}
		ctx.JSON(http.StatusUnauthorized, result)
		return
	}

	logResponse, err := controller.CustomerService.LogActivity(customerId)
	helper.PanicError(err)

	log.Info("View history customer id: " + customerId)
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "View history activity with ID customer: " + customerId,
		Data:    logResponse,
	}

	ctx.JSON(http.StatusOK, gin.H{"customer": webResponse})
}

// Edit implements CustomerController.
func (cc *CustomerControllerImpl) Edit(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(jwt.MapClaims)
	customerId := claims["id"].(string)
	role := claims["role"].(string)

	if role != "customer" {
		result := web.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "unauthorized",
			Data:    "user is unauthorized",
		}
		ctx.JSON(http.StatusUnauthorized, result)
		return
	}

	var req customerweb.UpdateRequest
	req.Id = customerId
	ctx.ShouldBindJSON(&req)

	response, err := cc.CustomerService.Edit(&req)
	helper.PanicError(err)

	log.Info("userId: " + customerId + " Has edit data profile")
	webResponse := web.WebResponse{
		Code:    http.StatusCreated,
		Status:  "CREATED",
		Message: "Update Customer Has Successfully",
		Data:    response,
	}

	ctx.JSON(http.StatusCreated, gin.H{"customer": webResponse})
}

// GetAll implements CustomerController.
func (cc *CustomerControllerImpl) GetAll(ctx *gin.Context) {
	response, err := cc.CustomerService.GetAll()
	helper.PanicError(err)

	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "List of all customer data",
		Data:    response,
	}

	ctx.JSON(http.StatusOK, gin.H{"customer": webResponse})
}

// GetByParams implements CustomerController.
func (cc *CustomerControllerImpl) GetByParams(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(jwt.MapClaims)
	customerId := claims["id"].(string)
	role := claims["role"].(string)

	if role != "customer" {
		result := web.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "unauthorized",
			Data:    "user is unauthorized",
		}
		ctx.JSON(http.StatusUnauthorized, result)
		return
	}

	response, _, err := cc.CustomerService.GetByParams(customerId, "")
	helper.PanicError(err)

	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Customer data by ID",
		Data:    response,
	}

	log.Info("userId: " + customerId + " Has view the profile ")
	ctx.JSON(http.StatusCreated, gin.H{"customer": webResponse})
}

// Register implements CustomerController.
func (cc *CustomerControllerImpl) Register(ctx *gin.Context) {
	var req customerweb.CreateRequest
	ctx.ShouldBindJSON(&req)

	response, err := cc.CustomerService.Register(&req)
	helper.PanicError(err)

	webResponse := web.WebResponse{
		Code:    http.StatusCreated,
		Status:  "CREATED",
		Message: "Create Customer Has Successfully",
		Data:    response,
	}
	log.Info("userId: " + response.Id + " Has create new customer account ")
	ctx.JSON(http.StatusCreated, gin.H{"customer": webResponse})
}

// Unreg implements CustomerController.
func (cc *CustomerControllerImpl) Unreg(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(jwt.MapClaims)
	customerId := claims["id"].(string)
	role := claims["role"].(string)

	if role != "customer" {
		result := web.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "unauthorized",
			Data:    "user is unauthorized",
		}
		ctx.JSON(http.StatusUnauthorized, result)
		return
	}
	response, err := cc.CustomerService.Unreg(customerId)
	helper.PanicError(err)

	log.Info("userId: " + customerId + " Has delete account ")
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Delete account has successfully",
		Data:    response,
	}

	ctx.JSON(http.StatusOK, gin.H{"customer": webResponse})
}

func NewCustomerController(customerService customersrv.CustomerService) CustomerController {
	return &CustomerControllerImpl{
		CustomerService: customerService,
	}
}
