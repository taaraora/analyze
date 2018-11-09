// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/supergiant/robot/pkg/models"
)

// GetCheckResultsOKCode is the HTTP code returned for type GetCheckResultsOK
const GetCheckResultsOKCode int = 200

/*GetCheckResultsOK no error

swagger:response getCheckResultsOK
*/
type GetCheckResultsOK struct {

	/*
	  In: Body
	*/
	Payload *GetCheckResultsOKBody `json:"body,omitempty"`
}

// NewGetCheckResultsOK creates GetCheckResultsOK with default headers values
func NewGetCheckResultsOK() *GetCheckResultsOK {

	return &GetCheckResultsOK{}
}

// WithPayload adds the payload to the get check results o k response
func (o *GetCheckResultsOK) WithPayload(payload *GetCheckResultsOKBody) *GetCheckResultsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get check results o k response
func (o *GetCheckResultsOK) SetPayload(payload *GetCheckResultsOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCheckResultsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetCheckResultsDefault error

swagger:response getCheckResultsDefault
*/
type GetCheckResultsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetCheckResultsDefault creates GetCheckResultsDefault with default headers values
func NewGetCheckResultsDefault(code int) *GetCheckResultsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetCheckResultsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get check results default response
func (o *GetCheckResultsDefault) WithStatusCode(code int) *GetCheckResultsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get check results default response
func (o *GetCheckResultsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get check results default response
func (o *GetCheckResultsDefault) WithPayload(payload *models.Error) *GetCheckResultsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get check results default response
func (o *GetCheckResultsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCheckResultsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}