package main

import (
	"context"
	"github.com/EvilKhaosKat/simple-peak-flowmeter/pkg/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

// RecordCtx middleware is used to load an Record object from
// the URL parameters passed through as the request. In case
// the Record could not be found, we stop here and return a 404.
func (app *application) RecordCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var Record *models.Record
		var err error

		if RecordID := chi.URLParam(r, "RecordID"); RecordID != "" {
			Record, err = app.records.Get(RecordID)
		} else {
			render.Render(w, r, ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "Record", Record)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}