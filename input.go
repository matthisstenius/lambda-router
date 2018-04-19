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

// PopulateBody with data in current request
func (i *Input) PopulateBody(out interface{}) error {
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

// PopulateEventBody parse body in SNS event
func (i *Input) PopulateEventBody(out interface{}) error {
	record := i.event["Records"].([]map[string]interface{})[0]
	message := record["SNS"].(map[string]string)["Message"]

	if err := json.Unmarshal([]byte(message), &out); err != nil {
		return errors.New("could not parse SNS Message")
	}
	return nil
}

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
