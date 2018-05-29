package http

import (
	"encoding/json"
	"errors"

	"github.com/matthisstenius/logger"
)

// Input for parsed HTTP event
type Input struct {
	event map[string]interface{}
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
		logger.WithFields(logger.Fields{
			"method": "HTTPInput",
			"body":   body,
		}).Error("missing request body")
		return errors.New("missing request body")
	}

	if err := json.Unmarshal([]byte(body.(string)), &out); err != nil {
		logger.WithFields(logger.Fields{
			"method": "HTTPInput",
			"error":  err,
		}).Error("could not parse body as JSON")
		return errors.New("could not parse body as JSON")
	}
	return nil
}

// Body un parsed body from HTTP event
func (i *Input) Body() string {
	return i.event["body"].(string)
}

// CurrentUser holding base information about currently authenticated user
type CurrentUser map[string]interface{}

// CurrentUser get currently authenticated user
func (i *Input) CurrentUser() CurrentUser {
	reqContext := i.event["requestContext"]

	authorizer, ok := reqContext.(map[string]interface{})["authorizer"]
	if !ok || authorizer == nil {
		logger.WithFields(logger.Fields{
			"method": "CurrentUser",
		}).Panic("authorizer index missing")
	}

	authData, ok := authorizer.(map[string]interface{})["authData"]
	if !ok || authorizer == nil {
		logger.WithFields(logger.Fields{
			"method": "CurrentUser",
		}).Panic("authorizer index missing")
	}

	var currentUser CurrentUser
	if value, ok := authData.(string); ok {
		err := json.Unmarshal([]byte(value), &currentUser)
		if err != nil {
			logger.WithFields(logger.Fields{
				"method": "CurrentUser",
				"error":  err,
			}).Panic("Could not parse authData")
		}
	} else {
		currentUser = authData.(map[string]interface{})
	}
	return currentUser
}
