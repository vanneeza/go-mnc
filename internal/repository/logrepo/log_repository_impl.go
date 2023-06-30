package logrepo

import (
	"database/sql"

	"github.com/vanneeza/go-mnc/internal/domain/entity"
	"github.com/vanneeza/go-mnc/utils/helper"
)

type LogRepositoryImpl struct {
}

// Create implements LogRepository.
func (*LogRepositoryImpl) Create(tx *sql.Tx, log entity.Log) (*entity.Log, error) {
	SQL := "INSERT INTO log_user(user_id, level, activity, description, created_at) VALUES($1, $2, $3, $4, $5)"
	_, err := tx.Exec(SQL, log.UserId, log.Level, log.Activity, log.Description, log.CreatedAt)
	helper.PanicError(err)

	return &log, nil
}

// GetAllByUserId implements LogRepository.
func (*LogRepositoryImpl) GetAllByUserId(tx *sql.Tx, userId string) (*[]entity.Log, error) {
	SQL := "SELECT id, user_id, level, activity, description, created_at FROM log_user WHERE user_id = $1"
	rows, err := tx.Query(SQL, userId)
	helper.PanicError(err)
	defer rows.Close()

	var logs []entity.Log
	for rows.Next() {
		var log entity.Log
		err2 := rows.Scan(&log.Id, &log.UserId, &log.Level, &log.Activity, &log.Description, &log.CreatedAt)
		helper.PanicError(err2)

		logs = append(logs, log)

	}

	return &logs, nil
}

func NewLogRepository() LogRepository {
	return &LogRepositoryImpl{}
}
