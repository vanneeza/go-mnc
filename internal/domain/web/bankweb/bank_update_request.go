package bankweb

import "github.com/vanneeza/go-mnc/internal/domain/web/merchantweb"

type UpdateRequest struct {
	Id            string               `json:"id"`
	Name          string               `json:"name"`
	BankAccount   string               `json:"bank_account"`
	Branch        string               `json:"branch"`
	AccountNumber int64                `json:"account_number"`
	Merchant      merchantweb.Response `json:"merchant"`
}

type BankAdminpdateRequest struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	BankAccount   string `json:"bank_account"`
	Branch        string `json:"branch"`
	AccountNumber int64  `json:"account_number"`
}
