package txctrl

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/vanneeza/go-mnc/internal/domain/web"
	"github.com/vanneeza/go-mnc/internal/domain/web/txweb"
	"github.com/vanneeza/go-mnc/internal/service/txsrv"
	"github.com/vanneeza/go-mnc/utils/helper"
	"github.com/vanneeza/go-mnc/utils/log"
	"github.com/vanneeza/go-mnc/utils/pkg"
)

type TxControllerImpl struct {
	TxService txsrv.TxService
}

// HistoryTransaction implements TxController.
func (*TxControllerImpl) HistoryTransaction(ctx *gin.Context) {
	panic("unimplemented")
}

// Confirmation implements TxController.
func (controller *TxControllerImpl) Confirmation(ctx *gin.Context) {
	detailId := ctx.Param("id")
	var req txweb.DetailUpdateRequest
	req.Id = detailId

	orderResponse, err := controller.TxService.Confirmation(&req)
	helper.PanicError(err)

	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "The payment order has confirmation by admin! Transaction Successfuly!",
		Data:    orderResponse,
	}
	ctx.JSON(http.StatusCreated, gin.H{"confirmation": webResponse})
}

// ViewAllOrder implements TxController.
func (controller *TxControllerImpl) ViewAllOrder(ctx *gin.Context) {
	orderResponse, err := controller.TxService.ViewAllPayment()
	helper.PanicError(err)

	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "List of all order data with status is DONE",
		Data:    orderResponse,
	}
	ctx.JSON(http.StatusCreated, gin.H{"orders": webResponse})
}

// Invoice implements TxController.
func (controller *TxControllerImpl) Invoice(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(jwt.MapClaims)
	customerId := claims["id"].(string)
	role := claims["role"].(string)

	if role != "customer" {
		result := web.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "user is unauthorized",
			Data:    "NULL",
		}
		ctx.JSON(http.StatusUnauthorized, result)
		return
	}

	var order txweb.OrderCreateRequest

	err := ctx.ShouldBindJSON(&order)
	if err != nil {
		log.Error("userId: "+customerId, err)
		webResponse := web.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "BAD_REQUEST",
			Message: "Failed to create order",
			Data:    "NULL",
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"invoice": webResponse})
		return
	}

	order.CustomerId = customerId
	orderResponse, err2 := controller.TxService.Invoice(&order)
	if err2 != nil {
		log.Error("userId: "+customerId, err2)
		webResponse := web.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "BAD_REQUEST",
			Message: err2.Error(),
			Data:    "NULL",
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"invoice": webResponse})
		return
	}

	log.Info("userId: " + customerId + " Has create the order")
	webResponse := web.WebResponse{
		Code:    http.StatusCreated,
		Status:  "CREATED",
		Message: "Create Invoice Successfully",
		Data:    orderResponse,
	}
	ctx.JSON(http.StatusCreated, gin.H{"invoice": webResponse})

}

func (controller *TxControllerImpl) Payment(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(jwt.MapClaims)
	customerId := claims["id"].(string)
	customerName := claims["name"].(string)
	role := claims["role"].(string)

	if role != "customer" {
		result := web.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "user is unauthorized",
			Data:    "NULL",
		}
		ctx.JSON(http.StatusUnauthorized, result)
		return
	}
	detailId := ctx.Param("id")

	var payment txweb.PaymentCreateRequest

	payment.CustomerId = customerId
	payment.DetailId = detailId

	err := ctx.ShouldBind(&payment)
	if err != nil {
		log.Error("userId: "+customerId, err)
		webResponse := web.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "BAD_REQUEST",
			Message: "Failed to create payment",
			Data:    "NULL",
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"invoice": webResponse})
		return
	}

	file, err2 := ctx.FormFile("photo")
	if err2 != nil {
		log.Error("userId: "+customerId, err)
		webResponse := web.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "BAD_REQUEST",
			Message: "Failed to create payment",
			Data:    "NULL",
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"invoice": webResponse})
		return
	}

	ext := filepath.Ext(file.Filename)
	validExtensions := []string{".png", ".jpg", ".jpeg"}

	isValidExtension := false
	for _, validExt := range validExtensions {
		if ext == validExt {
			isValidExtension = true
			break
		}
	}

	if !isValidExtension {
		webResponse := web.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "BAD_REQUEST",
			Message: "Invalid file extension. Only PNG, JPG, and JPEG are allowed.",
			Data:    "NULL",
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"invoice": webResponse})
		return
	}

	currentTime := time.Now()
	formattedDate := currentTime.Format("20060102")
	key := pkg.GenerateRandomNumber()

	newFilename := fmt.Sprintf("%s%s%s%s%s%s%d%s", formattedDate, "_", "CustomerPayment", "_", customerName, "_", key, ext)
	uploadPath := filepath.Join("utils/document/uploads/customer_payment", newFilename)
	err3 := ctx.SaveUploadedFile(file, uploadPath)
	helper.PanicError(err3)

	paymentResponse, errPayment := controller.TxService.Payment(&payment)
	if errPayment != nil {
		log.Error("userId: "+customerId, errPayment)
		webResponse := web.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "BAD_REQUEST",
			Message: errPayment.Error(),
			Data:    "NULL",
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"invoice": webResponse})
		return
	}

	log.Info("userId: " + customerId + " Has create the payment")
	webResponse := web.WebResponse{
		Code:    http.StatusCreated,
		Status:  "CREATED",
		Message: "the payment was done, waiting to confirmation!",
		Data:    paymentResponse,
	}

	ctx.JSON(http.StatusCreated, gin.H{"payment": webResponse})

}

func NewTxController(txService txsrv.TxService) TxController {
	return &TxControllerImpl{
		TxService: txService,
	}
}
