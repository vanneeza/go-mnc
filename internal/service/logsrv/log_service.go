package logsrv

import (
	"github.com/vanneeza/go-mnc/internal/domain/entity"
)

type LogService interface {
	Register(req *entity.Log)
	GetAll(userId string)
}
