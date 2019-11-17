package main

import (
	"github.com/EvilKhaosKat/simple-peak-flowmeter/pkg/models"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"net/http"
)

const ContextKeyRecord = "record"
const ContextKeyNewRecordValue = "newRecordValue"

// SimpleCreateRecord persists the Record and returns it
// back to the client as an acknowledgement.
func (app *application) SimpleCreateRecord(w http.ResponseWriter, r *http.Request) {
	newRecordValue := r.Context().Value(ContextKeyNewRecordValue).(float32)

	record := app.recordsService.NewRecordByValue(newRecordValue)

	_, err := app.records.Update(record.ID, record.CreatedAt, record.Value)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewRecordResponse(record))
}

// CreateRecord persists the Record and returns it
// back to the client as an acknowledgement.
func (app *application) CreateRecord(w http.ResponseWriter, r *http.Request) {
	data := &RecordRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	record := data.Record
	record.ID = uuid.New().String()
	app.records.Update(record.ID, record.CreatedAt, record.Value)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewRecordResponse(record))
}

// GetRecord returns the specific Record. You'll notice it just
// fetches the Record right off the context, as its understood that
// if we made it this far, the Record must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func (app *application) GetRecord(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the record
	// context because this handler is a child of the RecordCtx
	// middleware. The worst case, the recoverer middleware will save us.
	record := r.Context().Value(ContextKeyRecord).(*models.Record)

	if err := render.Render(w, r, NewRecordResponse(record)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func (app *application) ListRecords(w http.ResponseWriter, r *http.Request) {
	records, err := app.records.GetAll()
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	if err := render.RenderList(w, r, NewRecordListResponse(records)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// UpdateRecord updates an existing Record in our persistent store.
func (app *application) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	record := r.Context().Value(ContextKeyRecord).(*models.Record)

	data := &RecordRequest{Record: record}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	record = data.Record
	app.records.Update(record.ID, record.CreatedAt, record.Value)

	render.Render(w, r, NewRecordResponse(record))
}

// DeleteRecord removes an existing Record from our persistent store.
func (app *application) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	var err error

	// Assume if we've reach this far, we can access the record
	// context because this handler is a child of the RecordCtx
	// middleware. The worst case, the recoverer middleware will save us.
	record := r.Context().Value(ContextKeyRecord).(*models.Record)

	_, err = app.records.Remove(record.ID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewRecordResponse(record))
}
