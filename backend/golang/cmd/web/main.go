package main

import (
	"context"
	"flag"
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
	generateRoutesDoc *bool
}

var timeoutCtx, _ = context.WithTimeout(context.Background(), 7*time.Second)

func main() {
	routes := flag.Bool("routes", false, "Generate router documentation")
	addr := flag.String("addr", ":3333", "HTTP network address")
	dsn := flag.String("dsn", "mongodb://localhost:27017", "MongoDB data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Println("Connecting to MongoDB")
	client, err := mongodb.OpenDB(*dsn)
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
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting HTTP server on %s", *addr)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}
