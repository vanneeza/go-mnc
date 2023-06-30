package txweb

import (
	"github.com/vanneeza/go-mnc/internal/domain/web/bankweb"
	"github.com/vanneeza/go-mnc/internal/domain/web/customerweb"
	"github.com/vanneeza/go-mnc/internal/domain/web/merchantweb"
	"github.com/vanneeza/go-mnc/internal/domain/web/productweb"
)

type DetailResponse struct {
	Id         string                      `json:"id"`
	Status     string                      `json:"status"`
	TotalPrice float64                     `json:"total_price"`
	Bank       []bankweb.ResponseForDetail `json:"bank"`
	Photo      string                      `json:"photo"`
}

type OrderResponse struct {
	Id       string               `json:"id"`
	Qty      int8                 `json:"qty"`
	Product  productweb.Response  `json:"product"`
	Customer customerweb.Response `json:"customer"`
	Detail   DetailResponse       `json:"detail"`
}

type OrderResponseWithoutDetail struct {
	Id       string               `json:"id"`
	Qty      int8                 `json:"qty"`
	Product  productweb.Response  `json:"product"`
	Customer customerweb.Response `json:"customer"`
}

type PaymentOrder struct {
	Product  productweb.Response       `json:"product"`
	Customer customerweb.Response      `json:"customer"`
	Bank     bankweb.ResponseForDetail `json:"bank"`
}

type PaymentResponse struct {
	Id           string        `json:"id"`
	PaymentOrder PaymentOrder  `json:"order"`
	Pay          float64       `json:"pay"`
	Detail       DetailPayment `json:"detail"`
}

type Payout struct {
	Payout       float64                   `json:"payout"`
	BankMerchant bankweb.ResponseForDetail `json:"bank_merchant"`
}

type OrderDetail struct {
	Id         string                     `json:"id"`
	Status     string                     `json:"status"`
	TotalPrice float64                    `json:"total_price"`
	Pay        float64                    `json:"pay"`
	Bank       bankweb.ResponseForDetail  `json:"bank"`
	Photo      string                     `json:"photo"`
	Order      OrderResponseWithoutDetail `json:"order"`
}

type Confirmation struct {
	Id         string                     `json:"id"`
	Status     string                     `json:"status"`
	TotalPrice float64                    `json:"total_price"`
	Pay        float64                    `json:"pay"`
	Bank       bankweb.ResponseForDetail  `json:"bank"`
	Photo      string                     `json:"photo"`
	Order      OrderResponseWithoutDetail `json:"order"`
	Payout     Payout                     `json:"payout"`
}

type DetailPayment struct {
	TotalPrice float64 `json:"total_price"`
	Photo      string  `json:"photo"`
}

type TxDetailMerchant struct {
	Id         string  `json:"id"`
	Status     string  `json:"status"`
	TotalPrice float64 `json:"total_price"`
	Photo      string  `json:"photo"`
}

type TxHistoryMerchantResponse struct {
	Detail   TxDetailMerchant           `json:"detail"`
	Product  productweb.ProductMerchant `json:"product"`
	Customer customerweb.Response       `json:"customer"`
}

type TxHistoryCustomerResponse struct {
	Detail   TxDetailMerchant           `json:"detail"`
	Product  productweb.ProductMerchant `json:"product"`
	Merchant merchantweb.Response       `json:"merchant"`
}
