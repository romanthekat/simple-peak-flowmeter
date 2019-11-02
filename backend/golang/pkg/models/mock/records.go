package mock

import (
	"github.com/EvilKhaosKat/simple-peak-flowmeter/pkg/models"
	"time"
)

type RecordModel struct{}

func NewRecordsModel() *RecordModel {
	return &RecordModel{}
}

func (r RecordModel) Update(id string, value float32) (string, error) {
	panic("implement me")
}

func (r RecordModel) Get(id string) (*models.Record, error) {
	panic("implement me")
}

func (r RecordModel) Remove(id string) (int64, error) {
	panic("implement me")
}

func (r RecordModel) GetAll() ([]*models.Record, error) {
	return records, nil
}

// Records fixture data
var records = []*models.Record{
	{ID: "0", CreatedAt: time.Now().Add(-1 * (time.Hour * 72)), Value: 490},
	{ID: "1", CreatedAt: time.Now().Add(-1 * (time.Hour * 48)), Value: 505},
	{ID: "2", CreatedAt: time.Now().Add(-1 * (time.Hour * 44)), Value: 480},
	{ID: "3", CreatedAt: time.Now().Add(-1 * (time.Hour * 24)), Value: 525},
	{ID: "4", CreatedAt: time.Now().Add(-1 * (time.Hour * 20)), Value: 495},
	{ID: "5", CreatedAt: time.Now(), Value: 520},
}
