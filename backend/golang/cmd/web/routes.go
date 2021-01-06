package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/docgen"
	"github.com/go-chi/render"
	"net/http"
	"strings"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	handleCors(r)

	// RESTy routes for "Records" resource
	r.Route("/records", func(r chi.Router) {
		r.Get("/", app.ListRecords)

		r.Post("/", app.CreateRecord) // POST /Records

		r.Route("/simple-add/{NewRecordValue}", func(r chi.Router) {
			r.Use(app.RecordNewValueCtx)
			r.Get("/", app.SimpleCreateRecord)
		})

		r.Route("/{RecordID}", func(r chi.Router) {
			r.Use(app.RecordCtx)            // Load the *Record on the request context
			r.Get("/", app.GetRecord)       // GET /Records/123
			r.Put("/", app.UpdateRecord)    // PUT /Records/123
			r.Delete("/", app.DeleteRecord) // DELETE /Records/123
		})
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/", handleMimeType(app, fileServer))
	r.Handle("/static/", http.StripPrefix("/static", handleMimeType(app, fileServer)))

	app.handleRoutesFileGeneration(r)

	return r
}

func handleMimeType(app *application, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}*/

		path := r.URL.Path

		app.infoLog.Println("fileserver path:", path)

		var contentType string
		if strings.HasSuffix(path, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(path, ".js") {
			contentType = "text/javascript"
		} else if strings.HasSuffix(path, ".html") {
			contentType = "text/html"
		} else {
			contentType = "text/html"
		}

		w.Header().Add("Content-Type", contentType)

		next.ServeHTTP(w, r)
	})
}

func handleCors(r *chi.Mux) {
	corsSettings := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r.Use(corsSettings.Handler)
}

func (app *application) handleRoutesFileGeneration(r *chi.Mux) {
	// Passing -routes to the program will generate docs for the above
	// router definition. See the `routes.json` file in this folder for
	// the output.
	if app.generateRoutesDoc {
		// fmt.Println(docgen.JSONRoutesDoc(r))
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "github.com/EvilKhaosKat/simple-peak-flowmeter/backend/golang",
			Intro:       "Simple peak flowmeter golang generated docs.",
		}))
	}
}
