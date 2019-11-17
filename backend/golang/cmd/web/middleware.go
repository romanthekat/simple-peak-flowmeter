package main

import (
	"context"
	"github.com/EvilKhaosKat/simple-peak-flowmeter/pkg/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

// RecordCtx middleware is used to load an Record object from
// the URL parameters passed through as the request. In case
// the Record could not be found, we stop here and return a 404.
func (app *application) RecordCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var record *models.Record
		var err error

		if RecordID := chi.URLParam(r, "RecordID"); RecordID != "" {
			record, err = app.records.Get(RecordID)
		} else {
			render.Render(w, r, ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyRecord, record)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RecordNewValueCtx middleware is used to load an record value object from
// the URL parameters passed through as the request. In case of error returns 400
func (app *application) RecordNewValueCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var value float64
		var err error

		if valueStr := chi.URLParam(r, "NewRecordValue"); valueStr != "" {
			value, err = strconv.ParseFloat(valueStr, 32)
		} else {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyNewRecordValue, float32(value))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
