package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/render"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

func main() {
	flag.Parse()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// RESTy routes for "Records" resource
	r.Route("/records", func(r chi.Router) {
		r.Get("/", ListRecords)
		r.Post("/", CreateRecord) // POST /Records

		r.Route("/{RecordID}", func(r chi.Router) {
			r.Use(RecordCtx)            // Load the *Record on the request context
			r.Get("/", GetRecord)       // GET /Records/123
			r.Put("/", UpdateRecord)    // PUT /Records/123
			r.Delete("/", DeleteRecord) // DELETE /Records/123
		})
	})

	// Passing -routes to the program will generate docs for the above
	// router definition. See the `routes.json` file in this folder for
	// the output.
	if *routes {
		// fmt.Println(docgen.JSONRoutesDoc(r))
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "github.com/EvilKhaosKat/simple-peak-flowmeter/backend/golang",
			Intro:       "Simple peak flowmeter golang generated docs.",
		}))
		return
	}

	http.ListenAndServe(":3333", r)
}

func ListRecords(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, NewRecordListResponse(Records)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// RecordCtx middleware is used to load an Record object from
// the URL parameters passed through as the request. In case
// the Record could not be found, we stop here and return a 404.
func RecordCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var Record *Record
		var err error

		if RecordID := chi.URLParam(r, "RecordID"); RecordID != "" {
			Record, err = dbGetRecord(RecordID)
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

// CreateRecord persists the Record and returns it
// back to the client as an acknowledgement.
func CreateRecord(w http.ResponseWriter, r *http.Request) {
	data := &RecordRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	Record := data.Record
	dbNewRecord(Record)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewRecordResponse(Record))
}

// GetRecord returns the specific Record. You'll notice it just
// fetches the Record right off the context, as its understood that
// if we made it this far, the Record must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func GetRecord(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the Record
	// context because this handler is a child of the RecordCtx
	// middleware. The worst case, the recoverer middleware will save us.
	Record := r.Context().Value("Record").(*Record)

	if err := render.Render(w, r, NewRecordResponse(Record)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// UpdateRecord updates an existing Record in our persistent store.
func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	Record := r.Context().Value("Record").(*Record)

	data := &RecordRequest{Record: Record}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	Record = data.Record
	dbUpdateRecord(Record.ID, Record)

	render.Render(w, r, NewRecordResponse(Record))
}

// DeleteRecord removes an existing Record from our persistent store.
func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	var err error

	// Assume if we've reach this far, we can access the Record
	// context because this handler is a child of the RecordCtx
	// middleware. The worst case, the recoverer middleware will save us.
	Record := r.Context().Value("Record").(*Record)

	Record, err = dbRemoveRecord(Record.ID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewRecordResponse(Record))
}

//--
// Request and Response payloads for the REST api.
//
// The payloads embed the data model objects an
//
// In a real-world project, it would make sense to put these payloads
// in another file, or another sub-package.
//--

// RecordRequest is the request payload for Record data model.
//
// NOTE: It's good practice to have well defined request and response payloads
// so you can manage the specific inputs and outputs for clients, and also gives
// you the opportunity to transform data on input or output, for example
// on request, we'd like to protect certain fields and on output perhaps
// we'd like to include a computed field based on other values that aren't
// in the data model. Also, check out this awesome blog post on struct composition:
// http://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
type RecordRequest struct {
	*Record

	ProtectedID string `json:"id"` // override 'id' json to have more control
}

func (a *RecordRequest) Bind(r *http.Request) error {
	// a.Record is nil if no Record fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if a.Record == nil {
		return errors.New("missing required Record fields")
	}

	// a.User is nil if no Userpayload fields are sent in the request. In this app
	// this won't cause a panic, but checks in this Bind method may be required if
	// a.User or futher nested fields like a.User.Name are accessed elsewhere.

	// just a post-process after a decode..
	a.ProtectedID = "" // unset the protected ID
	return nil
}

// RecordResponse is the response payload for the Record data model.
// See NOTE above in RecordRequest as well.
//
// In the RecordResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type RecordResponse struct {
	*Record
}

func NewRecordResponse(Record *Record) *RecordResponse {
	resp := &RecordResponse{Record: Record}

	return resp
}

func (rd *RecordResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func NewRecordListResponse(Records []*Record) []render.Renderer {
	list := []render.Renderer{}
	for _, Record := range Records {
		list = append(list, NewRecordResponse(Record))
	}
	return list
}

// NOTE: as a thought, the request and response payloads for an Record could be the
// same payload type, perhaps will do an example with it as well.
// type RecordPayload struct {
//   *Record
// }

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found"}

//--
// Data model objects and persistence mocks:
//--

// User data model
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Record data model. I suggest looking at https://upper.io for an easy
// and powerful data persistence adapter.
type Record struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Value     float32   `json:"value"`
}

// Record fixture data
var Records = []*Record{
	{ID: "1", CreatedAt: time.Now().Add(-1 * (time.Hour * 48)), Value: 505},
	{ID: "2", CreatedAt: time.Now().Add(-1 * (time.Hour * 44)), Value: 480},
	{ID: "3", CreatedAt: time.Now().Add(-1 * (time.Hour * 24)), Value: 525},
	{ID: "4", CreatedAt: time.Now().Add(-1 * (time.Hour * 20)), Value: 495},
	{ID: "5", CreatedAt: time.Now(), Value: 520},
}

func dbNewRecord(Record *Record) (string, error) {
	Record.ID = fmt.Sprintf("%d", rand.Intn(100)+10)
	Records = append(Records, Record)
	return Record.ID, nil
}

func dbGetRecord(id string) (*Record, error) {
	for _, a := range Records {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, errors.New("record not found")
}

func dbUpdateRecord(id string, Record *Record) (*Record, error) {
	for i, a := range Records {
		if a.ID == id {
			Records[i] = Record
			return Record, nil
		}
	}
	return nil, errors.New("record not found")
}

func dbRemoveRecord(id string) (*Record, error) {
	for i, a := range Records {
		if a.ID == id {
			Records = append((Records)[:i], (Records)[i+1:]...)
			return a, nil
		}
	}
	return nil, errors.New("record not found")
}
