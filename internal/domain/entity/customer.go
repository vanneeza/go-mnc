package entity

type Customer struct {
	Id       string
	Name     string
	Phone    string
	Address  string
	Password string
	Role     string
}

type CustomerTxHistory struct {
	Detail   Detail
	Qty      string
	Product  Product
	Merchant Merchant
}
