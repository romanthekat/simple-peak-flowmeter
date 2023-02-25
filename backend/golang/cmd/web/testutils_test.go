package main

import (
	"github.com/romanthekat/simple-peak-flowmeter/pkg/models/mock"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func newTestApplication(t *testing.T) *application {
	recordsModel := mock.NewRecordsModel()
	return &application{
		errorLog:          log.New(ioutil.Discard, "", 0),
		infoLog:           log.New(ioutil.Discard, "", 0),
		records:           recordsModel,
		generateRoutesDoc: false,
	}
}

func newGetRequest(t *testing.T, url string) *http.Request {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	return r
}
