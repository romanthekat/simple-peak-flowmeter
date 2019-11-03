package mongodb

import (
	"context"
	"fmt"
	"github.com/EvilKhaosKat/simple-peak-flowmeter/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"unicode/utf8"
)

const (
	databaseName      = "simple-peak-flowmeter"
	collectionRecords = "records"
)

var (
	ctx = context.Background()
)

func OpenDB(dsn string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}

type RecordModel struct {
	client *mongo.Client
}

func NewRecordModel(client *mongo.Client) *RecordModel {
	return &RecordModel{client}
}

func (m *RecordModel) getRecordsCollection() *mongo.Collection {
	return m.client.Database(databaseName).Collection(collectionRecords)
}

// This will insert a new person into the database or updates existing.
func (m *RecordModel) Update(id string, value float32) (string, error) {
	persons := m.getRecordsCollection()

	upsert := true
	result, err := persons.UpdateOne(ctx,
		bson.M{"id": id},
		bson.M{
			"$set": bson.M{
				"id":    id,
				"value": value},
		},
		&options.UpdateOptions{
			Upsert: &upsert,
		},
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(result.UpsertedID), nil
}

// This will return a specific person based on its id.
func (m *RecordModel) Get(id string) (*models.Record, error) {
	if utf8.RuneCountInString(id) == 0 {
		return nil, nil
	}

	persons := m.getRecordsCollection()

	result := persons.FindOne(ctx, bson.M{"id": id})

	var record *models.Record
	err := result.Decode(&record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (m *RecordModel) Remove(id string) (int64, error) {
	if utf8.RuneCountInString(id) == 0 {
		return 0, nil
	}

	persons := m.getRecordsCollection()

	result, err := persons.DeleteOne(ctx, bson.M{"id": id})
	return result.DeletedCount, err
}

// This will return all the created persons.
func (m *RecordModel) GetAll() ([]*models.Record, error) {
	var result []*models.Record

	persons := m.getRecordsCollection()
	cur, err := persons.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var person models.Record
		err := cur.Decode(&person)
		if err != nil {
			return nil, err
		}

		result = append(result, &person)
	}
	return result, nil
}
