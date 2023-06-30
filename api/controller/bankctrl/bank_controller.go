package bankctrl

import "github.com/gin-gonic/gin"

type BankController interface {
	Register(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetByParams(ctx *gin.Context)
	Edit(ctx *gin.Context)
	Unreg(ctx *gin.Context)

	RegisterBankAdmin(ctx *gin.Context)
	GetAllBankAdmin(ctx *gin.Context)
	GetByIdBankAdmin(ctx *gin.Context)
}
