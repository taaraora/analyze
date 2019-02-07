// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// CheckResult CheckResult represents the single result of Check function invocation of specific plugin.
// swagger:model checkResult
type CheckResult struct {

	// shows check status
	// Enum: [RED YELLOW GREEN]
	CheckStatus string `json:"checkStatus,omitempty"`

	// date/Time of check execution
	// Format: date-time
	CompletedAt strfmt.DateTime `json:"completedAt,omitempty"`

	// detailed check result description, it basically contains plugin specific check result info
	Description interface{} `json:"description,omitempty"`

	// shows check execution errors
	ExecutionStatus string `json:"executionStatus,omitempty"`

	// unique UUID of Check function invocation of specific plugin
	ID string `json:"id,omitempty"`

	// check name
	Name string `json:"name,omitempty"`
}

// Validate validates this check result
func (m *CheckResult) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCheckStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCompletedAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var checkResultTypeCheckStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["RED","YELLOW","GREEN"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		checkResultTypeCheckStatusPropEnum = append(checkResultTypeCheckStatusPropEnum, v)
	}
}

const (

	// CheckResultCheckStatusRED captures enum value "RED"
	CheckResultCheckStatusRED string = "RED"

	// CheckResultCheckStatusYELLOW captures enum value "YELLOW"
	CheckResultCheckStatusYELLOW string = "YELLOW"

	// CheckResultCheckStatusGREEN captures enum value "GREEN"
	CheckResultCheckStatusGREEN string = "GREEN"
)

// prop value enum
func (m *CheckResult) validateCheckStatusEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, checkResultTypeCheckStatusPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *CheckResult) validateCheckStatus(formats strfmt.Registry) error {

	if swag.IsZero(m.CheckStatus) { // not required
		return nil
	}

	// value enum
	if err := m.validateCheckStatusEnum("checkStatus", "body", m.CheckStatus); err != nil {
		return err
	}

	return nil
}

func (m *CheckResult) validateCompletedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.CompletedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("completedAt", "body", "date-time", m.CompletedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *CheckResult) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CheckResult) UnmarshalBinary(b []byte) error {
	var res CheckResult
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
