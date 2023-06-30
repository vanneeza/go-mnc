package productweb

import "github.com/vanneeza/go-mnc/internal/domain/web/merchantweb"

type Response struct {
	Id          string               `json:"id"`
	Name        string               `json:"name"`
	Price       float64              `json:"price"`
	Description string               `json:"description"`
	Merchant    merchantweb.Response `json:"merchant"`
}

type ProductMerchant struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}
