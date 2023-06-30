package entity

type Bank struct {
	Id            string
	Name          string
	BankAccount   string
	Branch        string
	AccountNumber int64
	Merchant      Merchant
}

type BankAdmin struct {
	Id            string
	Name          string
	BankAccount   string
	Branch        string
	AccountNumber int64
}
