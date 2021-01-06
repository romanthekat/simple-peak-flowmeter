package main

import (
	"context"
	"github.com/EvilKhaosKat/simple-peak-flowmeter/pkg/models"
	"github.com/EvilKhaosKat/simple-peak-flowmeter/pkg/models/mongodb"
	"github.com/EvilKhaosKat/simple-peak-flowmeter/pkg/services"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	errorLog          *log.Logger
	infoLog           *log.Logger
	records           models.RecordModel
	recordsService    *services.RecordsService
	generateRoutesDoc bool
	authorizedIp      string
}

var timeoutCtx, _ = context.WithTimeout(context.Background(), 7*time.Second)

func main() {
	var routes bool
	if getEnv("ROUTES", "false") == "true" {
		routes = true
	}
	addr := getEnv("ADDR", ":3333")
	dsn := getEnv("DSN", "mongodb://mongo:27017")
	authorizedIp := getEnv("AUTHORIZED_IP", "-1")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Printf("Authorized IP: %s", authorizedIp)
	infoLog.Printf("DSN: %s", dsn)

	infoLog.Println("Connecting to MongoDB")
	client, err := mongodb.OpenDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer client.Disconnect(timeoutCtx)

	recordModel := mongodb.NewRecordModel(client)

	app := &application{
		errorLog:          errorLog,
		infoLog:           infoLog,
		records:           recordModel,
		recordsService:    services.NewRecordsService(),
		generateRoutesDoc: routes,
		authorizedIp:      authorizedIp,
	}

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting HTTP server on %s", addr)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
