package entity

type Product struct {
	Id          string
	Name        string
	Price       float64
	Description string
	Merchant    Merchant
}
