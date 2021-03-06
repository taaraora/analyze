// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/supergiant/analyze/pkg/models"
)

// ReplacePluginConfigOKCode is the HTTP code returned for type ReplacePluginConfigOK
const ReplacePluginConfigOKCode int = 200

/*ReplacePluginConfigOK plugin is removed from registry

swagger:response replacePluginConfigOK
*/
type ReplacePluginConfigOK struct {
}

// NewReplacePluginConfigOK creates ReplacePluginConfigOK with default headers values
func NewReplacePluginConfigOK() *ReplacePluginConfigOK {

	return &ReplacePluginConfigOK{}
}

// WriteResponse to the client
func (o *ReplacePluginConfigOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// ReplacePluginConfigNotFoundCode is the HTTP code returned for type ReplacePluginConfigNotFound
const ReplacePluginConfigNotFoundCode int = 404

/*ReplacePluginConfigNotFound Not Found

swagger:response replacePluginConfigNotFound
*/
type ReplacePluginConfigNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplacePluginConfigNotFound creates ReplacePluginConfigNotFound with default headers values
func NewReplacePluginConfigNotFound() *ReplacePluginConfigNotFound {

	return &ReplacePluginConfigNotFound{}
}

// WithPayload adds the payload to the replace plugin config not found response
func (o *ReplacePluginConfigNotFound) WithPayload(payload *models.Error) *ReplacePluginConfigNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace plugin config not found response
func (o *ReplacePluginConfigNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplacePluginConfigNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*ReplacePluginConfigDefault error

swagger:response replacePluginConfigDefault
*/
type ReplacePluginConfigDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplacePluginConfigDefault creates ReplacePluginConfigDefault with default headers values
func NewReplacePluginConfigDefault(code int) *ReplacePluginConfigDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplacePluginConfigDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace plugin config default response
func (o *ReplacePluginConfigDefault) WithStatusCode(code int) *ReplacePluginConfigDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace plugin config default response
func (o *ReplacePluginConfigDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the replace plugin config default response
func (o *ReplacePluginConfigDefault) WithPayload(payload *models.Error) *ReplacePluginConfigDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace plugin config default response
func (o *ReplacePluginConfigDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplacePluginConfigDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
