package productweb

import "github.com/vanneeza/go-mnc/internal/domain/web/merchantweb"

type CreateRequest struct {
	Name        string               `json:"name"`
	Price       float64              `json:"price"`
	Description string               `json:"description"`
	Merchant    merchantweb.Response `json:"merchant"`
}
