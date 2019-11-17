package services

import (
	"github.com/EvilKhaosKat/simple-peak-flowmeter/pkg/models"
	"time"
)
import "github.com/google/uuid"

type RecordsService struct {
}

func NewRecordsService() *RecordsService {
	return &RecordsService{}
}

func (r *RecordsService) NewRecordByValue(value float32) *models.Record {
	record := &models.Record{}

	record.ID = uuid.New().String()
	record.Value = value
	record.CreatedAt = time.Now()

	return record
}