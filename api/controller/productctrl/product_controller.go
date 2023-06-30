package productctrl

import "github.com/gin-gonic/gin"

type ProductController interface {
	Register(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetByParams(ctx *gin.Context)
	Edit(ctx *gin.Context)
	Unreg(ctx *gin.Context)
}
