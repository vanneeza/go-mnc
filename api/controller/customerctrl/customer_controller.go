package customerctrl

import "github.com/gin-gonic/gin"

type CustomerController interface {
	Register(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetByParams(ctx *gin.Context)
	Edit(ctx *gin.Context)
	Unreg(ctx *gin.Context)
	LogActivity(ctx *gin.Context)
	ViewTxHistories(ctx *gin.Context)
}
