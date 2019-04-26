// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/supergiant/analyze/pkg/models"
)

// NewReplacePluginConfigParams creates a new ReplacePluginConfigParams object
// no default values defined in spec.
func NewReplacePluginConfigParams() ReplacePluginConfigParams {

	return ReplacePluginConfigParams{}
}

// ReplacePluginConfigParams contains all the bound params for the replace plugin config operation
// typically these are obtained from a http.Request
//
// swagger:parameters replacePluginConfig
type ReplacePluginConfigParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: body
	*/
	Body *models.PluginConfig
	/*The id of the plugin to retrieve
	  Required: true
	  In: path
	*/
	PluginID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewReplacePluginConfigParams() beforehand.
func (o *ReplacePluginConfigParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.PluginConfig
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("body", "body"))
			} else {
				res = append(res, errors.NewParseError("body", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Body = &body
			}
		}
	} else {
		res = append(res, errors.Required("body", "body"))
	}
	rPluginID, rhkPluginID, _ := route.Params.GetOK("pluginId")
	if err := o.bindPluginID(rPluginID, rhkPluginID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindPluginID binds and validates parameter PluginID from path.
func (o *ReplacePluginConfigParams) bindPluginID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.PluginID = raw

	return nil
}
