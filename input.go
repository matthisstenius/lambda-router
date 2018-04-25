package api

import (
	"encoding/json"
	"errors"
	"reflect"

	"bitbucket.org/mstenius/logger"
	"github.com/asaskevich/govalidator"
)

// Input holding data for current request
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
			"method": "Input",
			"body":   body,
		}).Error("missing request body")
		return errors.New("missing request body")
	}

	if err := json.Unmarshal([]byte(body.(string)), &out); err != nil {
		logger.WithFields(logger.Fields{
			"method": "Input",
			"error":  err,
		}).Error("could not parse body as JSON")
		return errors.New("could not parse body as JSON")
	}

	t := reflect.ValueOf(out)
	if t.Elem().Kind() == reflect.Slice {
		t = t.Elem()
		for i := 0; i < t.Len(); i++ {
			if _, err := govalidator.ValidateStruct(t.Index(i)); err != nil {
				return err
			}
		}
		return nil
	}

	if _, err := govalidator.ValidateStruct(out); err != nil {
		return err
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
	err := json.Unmarshal([]byte(authData.(string)), &currentUser)
	if err != nil {
		logger.WithFields(logger.Fields{
			"method": "CurrentUser",
			"error":  err,
		}).Panic("Could not parse authData")
	}
	return currentUser
}

// StreamInput ...
type StreamInput struct {
	event map[string]interface{}
}

// ParseOldImage from DynamoDB stream event
func (si *StreamInput) ParseOldImage(out map[string]interface{}) error {
	record := si.event["Records"].([]interface{})[0]
	image, ok := record.(map[string]interface{})["dynamodb"].(map[string]interface{})["OldImage"].(map[string]interface{})
	if !ok {
		logger.WithFields(logger.Fields{
			"record": record,
		}).Error("StreamInput::ParseOldImage() missing OldImage attribute in event")
		return errors.New("missing OldImage attribute in event")
	}

	for k, attributes := range image {
		for _, attribute := range attributes.(map[string]interface{}) {
			out[k] = attribute
		}
	}
	return nil
}

// ParseNewImage from DynamoDB stream event
func (si *StreamInput) ParseNewImage(out map[string]interface{}) error {
	record := si.event["Records"].([]interface{})[0]
	image, ok := record.(map[string]interface{})["dynamodb"].(map[string]interface{})["NewImage"].(map[string]interface{})
	if !ok {
		logger.WithFields(logger.Fields{
			"record": record,
		}).Error("StreamInput::ParseNewImage() missing NewImage attribute in event")
		return errors.New("missing NewImage attribute in event")
	}

	for k, attributes := range image {
		for _, attribute := range attributes.(map[string]interface{}) {
			out[k] = attribute
		}
	}
	return nil
}
