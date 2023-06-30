package entity

type Merchant struct {
	Id       string
	Name     string
	Phone    string
	Password string
	Role     string
}

type MerchantTxHistory struct {
	Detail   Detail
	Qty      string
	Product  Product
	Customer Customer
}
