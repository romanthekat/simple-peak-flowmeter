package main

import (
	"github.com/EvilKhaosKat/simple-peak-flowmeter/pkg/models"
	"github.com/go-chi/render"
	"net/http"
)

// CreateRecord persists the Record and returns it
// back to the client as an acknowledgement.
func (app *application) CreateRecord(w http.ResponseWriter, r *http.Request) {
	data := &RecordRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	Record := data.Record
	app.records.Update(Record.ID, Record.Value)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewRecordResponse(Record))
}

// GetRecord returns the specific Record. You'll notice it just
// fetches the Record right off the context, as its understood that
// if we made it this far, the Record must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func (app *application) GetRecord(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the Record
	// context because this handler is a child of the RecordCtx
	// middleware. The worst case, the recoverer middleware will save us.
	Record := r.Context().Value("Record").(*models.Record)

	if err := render.Render(w, r, NewRecordResponse(Record)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func (app *application) ListRecords(w http.ResponseWriter, r *http.Request) {
	records, err := app.records.GetAll()
	if err != nil {
		render.Render(w,r, ErrRender(err))
		return
	}

	if err := render.RenderList(w, r, NewRecordListResponse(records)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// UpdateRecord updates an existing Record in our persistent store.
func (app *application) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	Record := r.Context().Value("Record").(*models.Record)

	data := &RecordRequest{Record: Record}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	Record = data.Record
	app.records.Update(Record.ID, Record.Value)

	render.Render(w, r, NewRecordResponse(Record))
}

// DeleteRecord removes an existing Record from our persistent store.
func (app *application) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	var err error

	// Assume if we've reach this far, we can access the Record
	// context because this handler is a child of the RecordCtx
	// middleware. The worst case, the recoverer middleware will save us.
	Record := r.Context().Value("Record").(*models.Record)

	_, err = app.records.Remove(Record.ID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewRecordResponse(Record))
}
