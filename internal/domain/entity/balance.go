package entity

type Balance struct {
	Id      string
	Balance float64
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}
