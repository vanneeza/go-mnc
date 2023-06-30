package entity

type Order struct {
	Id       string
	Qty      int8
	Product  Product
	Customer Customer
	Detail   Detail
}

type OrderDetail struct {
	Order  Order   `json:"order"`
	Detail Detail  `json:"detail"`
	Pay    float64 `json:"pay"`
}
