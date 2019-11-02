package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrDbProblem = errors.New("models: problem with db")

//Record struct contains information of one measurement record
type Record struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Value     float32   `json:"value"`
}

//RecordModel defines model/DAO methods for Record
type RecordModel interface {
	Update(id string, value float32) (string, error)
	Get(id string) (*Record, error)
	Remove(id string) (int64, error)
	GetAll() ([]*Record, error)
}
