package merchantctrl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/vanneeza/go-mnc/internal/domain/web"
	"github.com/vanneeza/go-mnc/internal/domain/web/merchantweb"
	"github.com/vanneeza/go-mnc/internal/service/merchantsrv"
	"github.com/vanneeza/go-mnc/utils/helper"
	"github.com/vanneeza/go-mnc/utils/log"
)

type MerchantControllerImpl struct {
	MerchantService merchantsrv.MerchantService
}

// HistoryTransaction implements MerchantController.
func (controller *MerchantControllerImpl) HistoryTransaction(ctx *gin.Context) {
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

	response, err := controller.MerchantService.MerchantTxHistory(merchantId)
	helper.PanicError(err)

	log.Info("View transaction history merchant id: " + merchantId)
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "List of all transaction history",
		Data:    response,
	}

	ctx.JSON(http.StatusOK, gin.H{"transaction_history": webResponse})

}

// LogActivity implements MerchantController.
func (controller *MerchantControllerImpl) LogActivity(ctx *gin.Context) {
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

	logResponse, err := controller.MerchantService.LogActivity(merchantId)
	helper.PanicError(err)

	log.Info("View Activity merchant id: " + merchantId)
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "View history activity with ID merchant: " + merchantId,
		Data:    logResponse,
	}

	ctx.JSON(http.StatusOK, gin.H{"merchant": webResponse})

}

// Edit implements MerchantController.
func (cc *MerchantControllerImpl) Edit(ctx *gin.Context) {
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

	var req merchantweb.UpdateRequest
	req.Id = merchantId
	errBindJson := ctx.ShouldBindJSON(&req)
	if errBindJson != nil {
		webResponse := web.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "INTERNAL SERVER ERROR",
			Message: "Sorry, please come back later!",
			Data:    "NULL",
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"merchant": webResponse})
		return
	}
	response, err := cc.MerchantService.Edit(&req)
	if err != nil {
		webResponse := web.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "BAD_REQUEST",
			Message: err.Error(),
			Data:    "NULL",
		}

		ctx.JSON(http.StatusCreated, gin.H{"merchant": webResponse})
		return
	}
	log.Info("Edit Merchant With MerchantID = " + merchantId)
	webResponse := web.WebResponse{
		Code:    http.StatusCreated,
		Status:  "CREATED",
		Message: "Update Merchant Has Successfully",
		Data:    response,
	}

	ctx.JSON(http.StatusCreated, gin.H{"merchant": webResponse})
}

// GetAll implements MerchantController.
func (cc *MerchantControllerImpl) GetAll(ctx *gin.Context) {

	response, err := cc.MerchantService.GetAll()
	helper.PanicError(err)

	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "List of all merchant data",
		Data:    response,
	}

	ctx.JSON(http.StatusCreated, gin.H{"merchant": webResponse})
}

// GetByParams implements MerchantController.
func (cc *MerchantControllerImpl) GetByParams(ctx *gin.Context) {
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

	response, _, err := cc.MerchantService.GetByParams(merchantId, "")
	helper.PanicError(err)

	log.Info(response.Name + " Has Viewed the Profile")
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Data Merchant With ID " + response.Id,
		Data:    response,
	}

	ctx.JSON(http.StatusCreated, gin.H{"merchant": webResponse})
}

// Register implements MerchantController.
func (cc *MerchantControllerImpl) Register(ctx *gin.Context) {
	var req merchantweb.CreateRequest
	ctx.ShouldBindJSON(&req)

	response, err := cc.MerchantService.Register(&req)
	helper.PanicError(err)

	webResponse := web.WebResponse{
		Code:    http.StatusCreated,
		Status:  "CREATED",
		Message: "Create Merchant Has Successfully",
		Data:    response,
	}

	log.Info("Register Merchant With Name " + req.Name)
	ctx.JSON(http.StatusCreated, gin.H{"merchant": webResponse})
}

// Unreg implements MerchantController.
func (cc *MerchantControllerImpl) Unreg(ctx *gin.Context) {
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
	response, err := cc.MerchantService.Unreg(merchantId)
	helper.PanicError(err)

	log.Info(response.Name + " Has Delete Account Merchant")
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Delete Account Merchant with ID " + merchantId,
		Data:    response,
	}

	ctx.JSON(http.StatusCreated, gin.H{"merchant": webResponse})
}

func NewMerchantController(merchantService merchantsrv.MerchantService) MerchantController {
	return &MerchantControllerImpl{
		MerchantService: merchantService,
	}
}
