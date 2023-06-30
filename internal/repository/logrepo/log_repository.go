package logrepo

import (
	"database/sql"

	"github.com/vanneeza/go-mnc/internal/domain/entity"
)

type LogRepository interface {
	Create(tx *sql.Tx, log entity.Log) (*entity.Log, error)
	GetAllByUserId(tx *sql.Tx, userId string) (*[]entity.Log, error)
}
