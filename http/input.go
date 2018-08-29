package http

import (
	"encoding/json"
	"errors"
	"github.com/matthisstenius/lambda-router/v2/domain"
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
	return value.(string)
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
	return value.(string)
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
