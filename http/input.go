package http

import (
	"encoding/json"
	"errors"
	"github.com/matthisstenius/lambda-router/v3/domain"
)

// Input for parsed HTTP event
type Input struct {
	event map[string]interface{}
}

// NewInput initializer
func NewInput(e map[string]interface{}) *Input {
	return &Input{event: e}
}

// GetPathParam in current request
func (i *Input) GetPathParam(param string) string {
	params, ok := i.event["pathParameters"]
	if !ok || params == nil {
		return ""
	}
	value, ok := params.(map[string]interface{})[param]
	if !ok {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	default:
		encoded, _ := json.Marshal(value)
		return string(encoded)
	}
}

// HasPathParam checks if param exists in path params
func (i *Input) HasPathParam(param string) bool {
	params, ok := i.event["pathParameters"]
	if !ok || params == nil {
		return false
	}
	if _, ok := params.(map[string]interface{})[param]; !ok {
		return false
	}
	return true
}

// HasQueryParam checks of param exists in request
func (i *Input) HasQueryParam(param string) bool {
	params, ok := i.event["queryStringParameters"]
	if !ok || params == nil {
		return false
	}

	if _, ok := params.(map[string]interface{})[param]; !ok {
		return false
	}
	return true
}

// GetQueryParam in current request
func (i *Input) GetQueryParam(param string) string {
	params, ok := i.event["queryStringParameters"]
	if !ok || params == nil {
		return ""
	}

	value, ok := params.(map[string]interface{})[param]
	if !ok {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	default:
		encoded, _ := json.Marshal(value)
		return string(encoded)
	}
}

// GetHeader in current request
func (i *Input) GetHeader(header string) string {
	params, ok := i.event["headers"]
	if !ok || params == nil {
		return ""
	}

	value, ok := params.(map[string]interface{})[header].(string)
	if !ok {
		return ""
	}
	return value
}

// ParseQueryParam in current request as JSON
func (i *Input) ParseQueryParam(param string, out interface{}) error {
	if err := json.Unmarshal([]byte(i.GetQueryParam(param)), &out); err != nil {
		return errors.New("could not parse param as JSON")
	}
	return nil
}

// ParseBody in current request
func (i *Input) ParseBody(out interface{}) error {
	body, ok := i.event["body"]
	if !ok || body == nil {
		return errors.New("missing request body")
	}

	if err := json.Unmarshal([]byte(body.(string)), &out); err != nil {
		return errors.New("could not parse body as JSON")
	}
	return nil
}

// Auth get auth properties based on given AuthProvider
func (i *Input) Auth(ap domain.AuthProvider) (domain.AuthProperties, error) {
	if ap == nil {
		return nil, errors.New("given auth provider is nil")
	}
	return ap.ParseAuth(i.event)
}

// RawBody get raw body form event
func (i *Input) RawBody() []byte {
	body, ok := i.event["body"]
	if !ok || body == nil {
		return []byte("")
	}
	return []byte(body.(string))
}
