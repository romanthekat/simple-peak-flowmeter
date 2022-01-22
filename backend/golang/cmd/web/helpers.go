package main

import (
	"errors"
	"github.com/romanthekat/simple-peak-flowmeter/pkg/models"
	"github.com/go-chi/render"
	"net/http"
)

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
	*models.Record

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
	*models.Record
}

func NewRecordResponse(Record *models.Record) *RecordResponse {
	resp := &RecordResponse{Record: Record}

	return resp
}

func (rd *RecordResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func NewRecordListResponse(Records []*models.Record) []render.Renderer {
	list := []render.Renderer{}
	for _, Record := range Records {
		list = append(list, NewRecordResponse(Record))
	}
	return list
}

func GetIPAddress(r *http.Request) string {
	return r.RemoteAddr
}
