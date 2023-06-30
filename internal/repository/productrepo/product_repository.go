package productrepo

import (
	"database/sql"

	"github.com/vanneeza/go-mnc/internal/domain/entity"
)

type ProductRepository interface {
	Create(tx *sql.Tx, product *entity.Product) (*entity.Product, error)
	FindAll(tx *sql.Tx) ([]entity.Product, error)
	FindByParams(tx *sql.Tx, productId, merchantId string) (*entity.Product, error)
	Update(tx *sql.Tx, product *entity.Product) (*entity.Product, error)
	Delete(tx *sql.Tx, productId string) error
}
