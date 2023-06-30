package bankweb

import "github.com/vanneeza/go-mnc/internal/domain/web/merchantweb"

type CreateRequest struct {
	Name          string               `json:"name"`
	BankAccount   string               `json:"bank_account"`
	Branch        string               `json:"branch"`
	AccountNumber int64                `json:"account_number"`
	Merchant      merchantweb.Response `json:"merchant"`
}

type BankAdminCreateRequest struct {
	Name          string `json:"name"`
	BankAccount   string `json:"bank_account"`
	Branch        string `json:"branch"`
	AccountNumber int64  `json:"account_number"`
}
