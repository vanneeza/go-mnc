package loginctrl

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/vanneeza/go-mnc/internal/domain/entity"
	"github.com/vanneeza/go-mnc/internal/domain/web"
	"github.com/vanneeza/go-mnc/internal/domain/web/loginweb"
	"github.com/vanneeza/go-mnc/internal/service/customersrv"
	"github.com/vanneeza/go-mnc/internal/service/logsrv"
	"github.com/vanneeza/go-mnc/internal/service/merchantsrv"
	"github.com/vanneeza/go-mnc/utils/helper"
	"github.com/vanneeza/go-mnc/utils/log"
)

type tokenData struct {
	Token string `json:"token"`
}

type LoginController interface {
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type loginController struct {
	customerService customersrv.CustomerService
	merchantService merchantsrv.MerchantService
	Validate        *validator.Validate
	LogService      logsrv.LogService
}

func NewLoginController(customerService customersrv.CustomerService, merchantService merchantsrv.MerchantService, validate *validator.Validate, logService logsrv.LogService) LoginController {
	return &loginController{
		customerService: customerService,
		merchantService: merchantService,
		Validate:        validate,
		LogService:      logService,
	}
}

func (l *loginController) Login(ctx *gin.Context) {
	var login loginweb.Request
	var role string
	jwtKey := os.Getenv("JWT_KEY")

	errBindJson := ctx.ShouldBindJSON(&login)
	if errBindJson != nil {
		webResponse := web.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "BAD_REQUEST",
			Message: "Failed to login",
			Data:    "NULL",
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"login": webResponse})
		return
	}

	err3 := l.Validate.Struct(login)
	if err3 != nil {
		webResponse := web.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "BAD_REQUEST",
			Message: "Failed to login",
			Data:    "NULL",
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"login": webResponse})
		return
	}

	_, customer, errCustomer := l.customerService.GetByParams("", login.Phone)
	_, merchant, errMerchant := l.merchantService.GetByParams("", login.Phone)

	if errCustomer != nil && errMerchant != nil {

		log.Error("failed to login", errors.New("user not found"))
		webResponse := web.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "BAD_REQUEST",
			Message: "Failed to login, user not found",
			Data:    "NULL",
		}
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	if customer != nil {
		role = "customer"
	}
	if merchant != nil {
		role = "merchant"
	}

	if role == "customer" {
		if customer.Phone == "" && customer.Password == "" {
			result := web.WebResponse{
				Code:    http.StatusNotFound,
				Status:  "NOT_FOUND",
				Message: "user not found",
				Data:    "NULL",
			}
			ctx.JSON(http.StatusNotFound, result)
			return
		} else {
			match := helper.CheckPasswordHash(login.Password, customer.Password)
			if !match {
				result := web.WebResponse{
					Code:    http.StatusBadRequest,
					Status:  "BAD_REQUEST",
					Message: "wrong password",
					Data:    "NULL",
				}
				ctx.JSON(http.StatusBadRequest, result)
				return
			}

			token := jwt.New(jwt.SigningMethodHS256)

			claims := token.Claims.(jwt.MapClaims)
			claims["id"] = customer.Id
			claims["name"] = customer.Name
			claims["role"] = customer.Role
			claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

			var jwtKeyByte = []byte(jwtKey)
			tokenString, err := token.SignedString(jwtKeyByte)
			helper.PanicError(err)
			ctx.SetCookie("Authorization", tokenString, 3600, "", "", false, true)
			log.Info(customer.Name + " Has Login")

			log := entity.Log{
				UserId:      customer.Id,
				Activity:    "Login",
				Description: customer.Name + " Has Login",
			}
			l.LogService.Register(&log)

			result := web.WebResponse{
				Code:    http.StatusOK,
				Status:  "OK",
				Message: "The customer has successfully logged in, hello " + customer.Name,
				Data:    "Token on cookies",
			}
			ctx.JSON(http.StatusOK, result)
		}
	}

	if role == "merchant" {
		if merchant.Phone == "" && merchant.Password == "" {
			result := web.WebResponse{
				Code:    http.StatusNotFound,
				Status:  "NOT_FOUND",
				Message: "user not found",
				Data:    "NULL",
			}
			ctx.JSON(http.StatusNotFound, result)
			return
		} else {
			match := helper.CheckPasswordHash(login.Password, merchant.Password)
			if !match {
				result := web.WebResponse{
					Code:    http.StatusBadRequest,
					Status:  "BAD_REQUEST",
					Message: "wrong password",
					Data:    "NULL",
				}
				ctx.JSON(http.StatusBadRequest, result)
				return
			}

			token := jwt.New(jwt.SigningMethodHS256)

			claims := token.Claims.(jwt.MapClaims)
			claims["id"] = merchant.Id
			claims["name"] = merchant.Name
			claims["phone"] = merchant.Phone
			claims["role"] = merchant.Role
			claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

			var jwtKeyByte = []byte(jwtKey)
			tokenString, err := token.SignedString(jwtKeyByte)
			helper.PanicError(err)

			ctx.SetCookie("Authorization", tokenString, 3600, "", "", false, true)
			log.Info(merchant.Name + " Has Login")

			log := entity.Log{
				UserId:      merchant.Id,
				Activity:    "Login",
				Description: merchant.Name + " Has Login",
			}
			l.LogService.Register(&log)

			result := web.WebResponse{
				Code:    http.StatusOK,
				Status:  "OK",
				Message: "The merchant has successfully logged in, hello " + merchant.Name,
				Data:    "Token on cookies",
			}
			ctx.JSON(http.StatusOK, result)
		}
	}

}

func (l *loginController) Logout(ctx *gin.Context) {

	claims := ctx.MustGet("claims").(jwt.MapClaims)
	userId := claims["id"].(string)
	name := claims["name"].(string)
	role := claims["role"].(string)

	if role == "merchant" {
		log.Info(name + " Has Logout")
		log := entity.Log{
			UserId:      userId,
			Activity:    "Logout",
			Description: name + " Has Logout",
		}
		l.LogService.Register(&log)
	} else {
		log.Info(name + " Has Logout")
		log := entity.Log{
			UserId:      userId,
			Activity:    "Logout",
			Description: name + " Has Logout",
		}
		l.LogService.Register(&log)
	}

	ctx.SetCookie("Authorization", "", -1, "", "", false, true)

	result := web.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Logout successfully, see you soon! " + name,
		Data:    "NULL",
	}
	ctx.JSON(http.StatusOK, result)

}
