package entity

import "time"

type Log struct {
	Id          int64     `json:"id"`
	UserId      string    `json:"user_id"`
	Level       string    `json:"level"`
	Activity    string    `json:"activity"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
