package main

import (
	"encoding/json"
	"github.com/EvilKhaosKat/simple-peak-flowmeter/pkg/models"
	"github.com/EvilKhaosKat/simple-peak-flowmeter/pkg/models/mock"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRecord(t *testing.T) {
	//given
	app := newTestApplication(t)

	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	r := newGetRequest(t, ts.URL+"/records/1")

	//when
	rs, err := ts.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}

	//then
	if rs.StatusCode != http.StatusOK {
		t.Fatalf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	var record *models.Record

	err = json.NewDecoder(rs.Body).Decode(&record)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%+v", mock.Records[1])
	if !isSameRecords(record, mock.Records[1]) {
		t.Errorf("want body contains json with mock record, got %+v", record)
	}

}
func TestGetAllRecord(t *testing.T) {
	//given
	app := newTestApplication(t)

	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	r := newGetRequest(t, ts.URL+"/records")

	//when
	rs, err := ts.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}

	//then
	if rs.StatusCode != http.StatusOK {
		t.Fatalf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	var Records []*models.Record

	err = json.NewDecoder(rs.Body).Decode(&Records)
	if err != nil {
		t.Fatal(err)
	}

	if len(Records) == 0 {
		t.Fatalf("Records must be provided, found none")
	}

	for index, record := range Records {
		if !isSameRecords(record, mock.Records[index]) {
			t.Errorf("want body contains json with mock Record, got %+v", record)
		}
	}
}

func isSameRecords(expected, actual *models.Record) bool {
	return expected.ID == actual.ID &&
		expected.Value == actual.Value &&
		expected.CreatedAt.Equal(actual.CreatedAt) //time.Time doesn't work with reflect.DeepEqual
}
