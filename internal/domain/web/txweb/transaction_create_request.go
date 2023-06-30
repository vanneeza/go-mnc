package txweb

import (
	"mime/multipart"
)

type OrderCreateRequest struct {
	Qty        int8   `json:"qty" validate:"required,max=100"`
	ProductId  string `json:"product" validate:"required"`
	CustomerId string `json:"customer"`
}

type DetailCreateRequest struct {
	Status         string  `json:"status"`
	TotalPrice     float64 `json:"total_price"`
	VirtualAccount int64   `json:"virtual_account"`
	Photo          string  `json:"photo"`
}

type PaymentCreateRequest struct {
	BankId     string                `form:"bank" validate:"required"`
	Pay        float64               `form:"pay" validate:"required,gt=0"`
	Photo      *multipart.FileHeader `form:"photo" validate:"required"`
	DetailId   string                `form:"detail" validate:"required"`
	CustomerId string                `form:"customer"`
}

type TestRequest struct {
	BankId     string                `form:"bank"`
	Pay        float64               `form:"pay"`
	Photo      *multipart.FileHeader `form:"photo"`
	DetailId   string                `form:"detail"`
	CustomerId string                `form:"customer"`
}
