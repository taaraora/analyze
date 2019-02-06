// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"./pkg/models"
)

// PatchPromethiusIntegrationInfoOKCode is the HTTP code returned for type PatchPromethiusIntegrationInfoOK
const PatchPromethiusIntegrationInfoOKCode int = 200

/*PatchPromethiusIntegrationInfoOK no error

swagger:response patchPromethiusIntegrationInfoOK
*/
type PatchPromethiusIntegrationInfoOK struct {

	/*
	  In: Body
	*/
	Payload *models.IntegrationInfo `json:"body,omitempty"`
}

// NewPatchPromethiusIntegrationInfoOK creates PatchPromethiusIntegrationInfoOK with default headers values
func NewPatchPromethiusIntegrationInfoOK() *PatchPromethiusIntegrationInfoOK {

	return &PatchPromethiusIntegrationInfoOK{}
}

// WithPayload adds the payload to the patch promethius integration info o k response
func (o *PatchPromethiusIntegrationInfoOK) WithPayload(payload *models.IntegrationInfo) *PatchPromethiusIntegrationInfoOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the patch promethius integration info o k response
func (o *PatchPromethiusIntegrationInfoOK) SetPayload(payload *models.IntegrationInfo) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PatchPromethiusIntegrationInfoOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PatchPromethiusIntegrationInfoDefault error

swagger:response patchPromethiusIntegrationInfoDefault
*/
type PatchPromethiusIntegrationInfoDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPatchPromethiusIntegrationInfoDefault creates PatchPromethiusIntegrationInfoDefault with default headers values
func NewPatchPromethiusIntegrationInfoDefault(code int) *PatchPromethiusIntegrationInfoDefault {
	if code <= 0 {
		code = 500
	}

	return &PatchPromethiusIntegrationInfoDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the patch promethius integration info default response
func (o *PatchPromethiusIntegrationInfoDefault) WithStatusCode(code int) *PatchPromethiusIntegrationInfoDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the patch promethius integration info default response
func (o *PatchPromethiusIntegrationInfoDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the patch promethius integration info default response
func (o *PatchPromethiusIntegrationInfoDefault) WithPayload(payload *models.Error) *PatchPromethiusIntegrationInfoDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the patch promethius integration info default response
func (o *PatchPromethiusIntegrationInfoDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PatchPromethiusIntegrationInfoDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}