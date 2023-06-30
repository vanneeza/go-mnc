package txsrv

import "github.com/vanneeza/go-mnc/internal/domain/web/txweb"

type TxService interface {
	Invoice(req *txweb.OrderCreateRequest) (*txweb.OrderResponse, error)
	Payment(req *txweb.PaymentCreateRequest) (*txweb.PaymentResponse, error)
	ViewAllPayment() ([]txweb.OrderDetail, error)
	Confirmation(req *txweb.DetailUpdateRequest) ([]txweb.Confirmation, error)
}
