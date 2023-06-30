package txctrl

import "github.com/gin-gonic/gin"

type TxController interface {
	Invoice(ctx *gin.Context)
	Payment(ctx *gin.Context)
	ViewAllOrder(ctx *gin.Context)
	Confirmation(ctx *gin.Context)
}
