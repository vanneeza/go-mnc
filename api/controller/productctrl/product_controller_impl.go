package productctrl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/vanneeza/go-mnc/internal/domain/web"
	"github.com/vanneeza/go-mnc/internal/domain/web/productweb"
	"github.com/vanneeza/go-mnc/internal/service/productsrv"
	"github.com/vanneeza/go-mnc/utils/helper"
	"github.com/vanneeza/go-mnc/utils/log"
)

type ProductControllerImpl struct {
	ProductService productsrv.ProductService
}

// Edit implements ProductController.
func (controller *ProductControllerImpl) Edit(ctx *gin.Context) {
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

	var req productweb.UpdateRequest
	productId := ctx.Param("id")
	req.Id = productId

	req.Merchant.Id = merchantId
	ctx.ShouldBindJSON(&req)

	response, err := controller.ProductService.Edit(&req)
	helper.PanicError(err)

	log.Info("userId: " + merchantId + " Has edit data product")
	webResponse := web.WebResponse{
		Code:    http.StatusCreated,
		Status:  "CREATED",
		Message: "Update Product Has Successfully",
		Data:    response,
	}

	ctx.JSON(http.StatusCreated, gin.H{"product": webResponse})
}

// GetAll implements ProductController.
func (controller *ProductControllerImpl) GetAll(ctx *gin.Context) {
	response, err := controller.ProductService.GetAll()
	helper.PanicError(err)

	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "List of all product data",
		Data:    response,
	}

	ctx.JSON(http.StatusOK, gin.H{"product": webResponse})
}

// GetByParams implements ProductController.
func (controller *ProductControllerImpl) GetByParams(ctx *gin.Context) {
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

	productId := ctx.Param("id")

	response, err := controller.ProductService.GetByParams(productId, "")
	helper.PanicError(err)

	log.Info("userId: " + merchantId + " Has view data product by ID")
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Product data from ID product " + productId,
		Data:    response,
	}

	ctx.JSON(http.StatusOK, gin.H{"product": webResponse})
}

// Register implements ProductController.
func (controller *ProductControllerImpl) Register(ctx *gin.Context) {
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

	var req productweb.CreateRequest
	ctx.ShouldBindJSON(&req)
	req.Merchant.Id = merchantId

	response, err := controller.ProductService.Register(&req)
	helper.PanicError(err)

	log.Info("userId: " + merchantId + " has add the new product")
	webResponse := web.WebResponse{
		Code:    http.StatusCreated,
		Status:  "CREATED",
		Message: "Created Product Has Successfully",
		Data:    response,
	}

	ctx.JSON(http.StatusCreated, gin.H{"product": webResponse})
}

// Unreg implements ProductController.
func (controller *ProductControllerImpl) Unreg(ctx *gin.Context) {
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

	productId := ctx.Param("id")
	response, err := controller.ProductService.Unreg(productId)
	helper.PanicError(err)

	log.Info("userId: " + merchantId + " has delete the product")
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Delete product data with ID " + productId + " Successfully!",
		Data:    response,
	}

	ctx.JSON(http.StatusCreated, gin.H{"product": webResponse})
}

func NewProductController(productService productsrv.ProductService) ProductController {
	return &ProductControllerImpl{
		ProductService: productService,
	}
}
