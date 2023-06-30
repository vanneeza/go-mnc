package txrepo

import (
	"database/sql"

	"github.com/vanneeza/go-mnc/internal/domain/entity"
)

type TxRepository interface {
	CreateOrder(tx *sql.Tx, order *entity.Order) (*entity.Order, error)
	CreateDetail(tx *sql.Tx, detail *entity.Detail) (*entity.Detail, error)
	CreatePayment(tx *sql.Tx, payment *entity.Payment) (*entity.Payment, error)
	UpdateDetail(tx *sql.Tx, detail *entity.Detail) (*entity.Detail, error)
	FindOrder(tx *sql.Tx, detailId string) (*entity.Order, error)
	GetAllOrder(tx *sql.Tx, detailId, status string) ([]entity.OrderDetail, error)
	ConfirmationOrder(tx *sql.Tx, detail *entity.Detail) (*entity.Detail, error)
}
