package bankctrl

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/vanneeza/go-mnc/internal/domain/web"
	"github.com/vanneeza/go-mnc/internal/domain/web/bankweb"
	"github.com/vanneeza/go-mnc/internal/service/banksrv"
	"github.com/vanneeza/go-mnc/utils/helper"
	"github.com/vanneeza/go-mnc/utils/log"
)

type BankControllerImpl struct {
	BankService banksrv.BankService
}

// GetByIdBankAdmin implements BankController.
func (controller *BankControllerImpl) GetByIdBankAdmin(ctx *gin.Context) {
	bankAdminId := ctx.Param("id")
	response, err := controller.BankService.GetByIdBankAdmin(bankAdminId)
	helper.PanicError(err)

	fmt.Printf("ctrl response: %v\n", response)
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "admin bank data with ID " + bankAdminId,
		Data:    response,
	}

	ctx.JSON(http.StatusCreated, gin.H{"admin/bank": webResponse})
}

// GetAllBankAdmin implements BankController.
func (controller *BankControllerImpl) GetAllBankAdmin(ctx *gin.Context) {
	response, err := controller.BankService.GetAllBankAdmin()
	helper.PanicError(err)

	fmt.Printf("ctrl response: %v\n", response)
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "List of all bank admin data",
		Data:    response,
	}

	ctx.JSON(http.StatusCreated, gin.H{"admin/bank": webResponse})
}

// RegisterBankAdmin implements BankController.
func (controller *BankControllerImpl) RegisterBankAdmin(ctx *gin.Context) {

	var req bankweb.BankAdminCreateRequest
	ctx.ShouldBindJSON(&req)

	fmt.Printf("ctrl req: %v\n", req)
	response, err := controller.BankService.RegisterBankAdmin(&req)
	helper.PanicError(err)

	webResponse := web.WebResponse{
		Code:    http.StatusCreated,
		Status:  "CREATED",
		Message: "Created Bank Has Successfully",
		Data:    response,
	}

	ctx.JSON(http.StatusCreated, gin.H{"admin/bank": webResponse})
}

// Edit implements BankController.
func (controller *BankControllerImpl) Edit(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(jwt.MapClaims)
	merchantId := claims["id"].(string)
	role := claims["role"].(string)

	if role != "merchant" {
		result := web.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "unauthorized",
			Data:    "user is unauthorized",
		}
		ctx.JSON(http.StatusUnauthorized, result)
		return
	}

	var req bankweb.UpdateRequest
	bankId := ctx.Param("id")
	req.Id = bankId

	req.Merchant.Id = merchantId
	ctx.ShouldBindJSON(&req)

	response, err := controller.BankService.Edit(&req)
	helper.PanicError(err)

	log.Info("userId: " + merchantId + " Has edit data bank")
	webResponse := web.WebResponse{
		Code:    http.StatusCreated,
		Status:  "CREATED",
		Message: "Update Bank Has Successfully",
		Data:    response,
	}

	ctx.JSON(http.StatusCreated, gin.H{"bank": webResponse})
}

// GetAll implements BankController.
func (controller *BankControllerImpl) GetAll(ctx *gin.Context) {
	response, err := controller.BankService.GetAll()
	helper.PanicError(err)

	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "List of all bank data",
		Data:    response,
	}

	ctx.JSON(http.StatusCreated, gin.H{"bank": webResponse})
}

// GetByParams implements BankController.
func (controller *BankControllerImpl) GetByParams(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(jwt.MapClaims)
	role := claims["role"].(string)
	merchantId := claims["id"].(string)

	if role != "merchant" {
		result := web.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "unauthorized",
			Data:    "user is unauthorized",
		}
		ctx.JSON(http.StatusUnauthorized, result)
		return
	}

	bankId := ctx.Param("id")

	response, err := controller.BankService.GetByParams(bankId, "")
	helper.PanicError(err)

	log.Info("userId: " + merchantId + " Has get data bank by ID")
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Bank data from ID bank " + bankId,
		Data:    response,
	}

	ctx.JSON(http.StatusCreated, gin.H{"bank": webResponse})
}

// Register implements BankController.
func (controller *BankControllerImpl) Register(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(jwt.MapClaims)
	merchantId := claims["id"].(string)
	role := claims["role"].(string)

	if role != "merchant" {
		result := web.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "unauthorized",
			Data:    "user is unauthorized",
		}
		ctx.JSON(http.StatusUnauthorized, result)
		return
	}

	var req bankweb.CreateRequest
	ctx.ShouldBindJSON(&req)
	req.Merchant.Id = merchantId

	response, err := controller.BankService.Register(&req)
	helper.PanicError(err)

	log.Info("userId: " + merchantId + " Has add new data bank")
	webResponse := web.WebResponse{
		Code:    http.StatusCreated,
		Status:  "CREATED",
		Message: "Created Bank Has Successfully",
		Data:    response,
	}

	ctx.JSON(http.StatusCreated, gin.H{"bank": webResponse})
}

// Unreg implements BankController.
func (controller *BankControllerImpl) Unreg(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(jwt.MapClaims)
	merchantId := claims["id"].(string)
	role := claims["role"].(string)

	if role != "merchant" {
		result := web.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "unauthorized",
			Data:    "user is unauthorized",
		}
		ctx.JSON(http.StatusUnauthorized, result)
		return
	}
	bankId := ctx.Param("id")

	response, err := controller.BankService.Unreg(bankId)
	helper.PanicError(err)

	log.Info("userId: " + merchantId + " Has add new data bank")
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Delete bank data with ID " + bankId + " Successfully!",
		Data:    response,
	}

	ctx.JSON(http.StatusCreated, gin.H{"bank": webResponse})
}

func NewBankController(bankService banksrv.BankService) BankController {
	return &BankControllerImpl{
		BankService: bankService,
	}
}
