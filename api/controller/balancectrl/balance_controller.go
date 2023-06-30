package balancectrl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vanneeza/go-mnc/internal/domain/web"
	"github.com/vanneeza/go-mnc/internal/service/balancesrv"
)

type BalanceController interface {
	ViewAll(ctx *gin.Context)
}

type BalanceControllerImpl struct {
	BalanceService balancesrv.BalanceService
}

// ViewAll implements BalanceController.
func (controller *BalanceControllerImpl) ViewAll(ctx *gin.Context) {
	balance, err := controller.BalanceService.ViewAll()
	if err != nil {
		webResponse := web.WebResponse{
			Code:    http.StatusNotFound,
			Status:  "STATUS_NOT_FOUND",
			Message: "balance not found",
			Data:    nil,
		}
		ctx.JSON(http.StatusNotFound, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Total Balance",
		Data:    balance,
	}
	ctx.JSON(http.StatusOK, webResponse)
}

func NewBalanceController(balanceService balancesrv.BalanceService) BalanceController {
	return &BalanceControllerImpl{
		BalanceService: balanceService,
	}
}
