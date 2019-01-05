package http

import (
	"encoding/json"
	"errors"
	"github.com/matthisstenius/lambda-router/v3/domain"
	"github.com/matthisstenius/logger"
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
func (i *Input) Auth() (*domain.AuthClaims, error) {
	reqContext := i.event["requestContext"]
	authorizer, ok := reqContext.(map[string]interface{})["authorizer"]
	if !ok || authorizer == nil {
		return nil, errors.New("authorizer index missing in event")
	}

	claims, ok := authorizer.(map[string]interface{})["claims"]
	if !ok {
		return nil, errors.New("claims index missing in authorizer")
	}

	var authProps map[string]interface{}
	if value, ok := claims.(string); ok {
		err := json.Unmarshal([]byte(value), &authProps)
		if err != nil {
			logger.WithFields(logger.Fields{
				"error": err,
			}).Error("CognitoAuthProvider::ParseAuth() Could not parse claims as JSON")
			return nil, errors.New("could not parse claims as JSON")
		}
	} else {
		authProps = claims.(map[string]interface{})
	}
	return domain.NewAuthClaims(authProps), nil
}

// RawBody get raw body form event
func (i *Input) RawBody() []byte {
	body, ok := i.event["body"]
	if !ok || body == nil {
		return []byte("")
	}
	return []byte(body.(string))
}
