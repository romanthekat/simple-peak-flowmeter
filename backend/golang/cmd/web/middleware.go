package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/romanthekat/simple-peak-flowmeter/pkg/models"
	"net/http"
	"strconv"
	"strings"
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

// RecordNewValueCtx middleware is used to load a record value object from
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

// LimitAuthorizedIp middleware limits requests to be performed from certain ip only
func (app *application) LimitAuthorizedIp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callerIp := GetIPAddress(r)

		if !strings.Contains(callerIp, app.authorizedIp) {
			app.errorLog.Printf("Authorized IP check failed, must be %s, request from %s\n",
				app.authorizedIp, callerIp)
			return
		}

		app.infoLog.Println("Authorized IP check passed")

		next.ServeHTTP(w, r)
	})
}
