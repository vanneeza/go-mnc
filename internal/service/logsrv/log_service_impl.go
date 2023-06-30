package logsrv

import (
	"database/sql"
	"time"

	"github.com/vanneeza/go-mnc/internal/domain/entity"
	"github.com/vanneeza/go-mnc/internal/repository/logrepo"
	"github.com/vanneeza/go-mnc/utils/helper"
)

type LogServiceImpl struct {
	Db            *sql.DB
	LogRepository logrepo.LogRepository
}

// GetAll implements LogService.
func (service LogServiceImpl) GetAll(userId string) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	service.LogRepository.GetAllByUserId(tx, userId)
}

// Register implements LogService.
func (service LogServiceImpl) Register(req *entity.Log) {
	tx, err := service.Db.Begin()
	helper.PanicError(err)
	defer helper.CommitOrRollback(tx)

	log := entity.Log{
		UserId:      req.UserId,
		Level:       "info",
		Activity:    req.Activity,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}
	service.LogRepository.Create(tx, log)
}

func NewLogService(Db *sql.DB, logRepository logrepo.LogRepository) LogService {
	return &LogServiceImpl{
		Db:            Db,
		LogRepository: logRepository,
	}
}
