package merchantctrl

import "github.com/gin-gonic/gin"

type MerchantController interface {
	Register(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetByParams(ctx *gin.Context)
	Edit(ctx *gin.Context)
	Unreg(ctx *gin.Context)
	LogActivity(ctx *gin.Context)
	HistoryTransaction(ctx *gin.Context)
}
