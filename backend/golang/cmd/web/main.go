package main

import (
	"flag"
	"github.com/EvilKhaosKat/simple-peak-flowmeter/pkg/models"
	"github.com/EvilKhaosKat/simple-peak-flowmeter/pkg/models/mock"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog          *log.Logger
	infoLog           *log.Logger
	records           models.RecordModel
	generateRoutesDoc *bool
}

func main() {
	routes := flag.Bool("routes", false, "Generate router documentation")
	addr := flag.String("addr", ":3333", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	recordModel := mock.NewRecordsModel()

	app := &application{
		errorLog:          errorLog,
		infoLog:           infoLog,
		records:           recordModel,
		generateRoutesDoc: routes,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting HTTP server on %s", *addr)
	err := srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}
