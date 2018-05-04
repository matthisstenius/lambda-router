package api

import (
	"encoding/json"
	"errors"
	"reflect"

	"strconv"

	"strings"

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

const (
	StreamEventInsert StreamEventType = "INSERT"
	StreamEventModify StreamEventType = "MODIFY"
	StreamEventRemove StreamEventType = "REMOVE"
)

// StreamInput ...
type StreamInput struct {
	event map[string]interface{}
}

// ParseOldImage from DynamoDB stream event
func (si *StreamInput) ParseOldImage(out interface{}) error {
	record := si.event["Records"].([]interface{})[0]
	image, ok := record.(map[string]interface{})["dynamodb"].(map[string]interface{})["OldImage"].(map[string]interface{})
	if !ok {
		logger.WithFields(logger.Fields{
			"record": record,
		}).Error("StreamInput::ParseOldImage() missing OldImage attribute in event")
		return errors.New("missing OldImage attribute in event")
	}

	if err := si.unmarshalAttributes(image, out); err != nil {
		return err
	}
	return nil
}

// ParseNewImage from DynamoDB stream event
func (si *StreamInput) ParseNewImage(out interface{}) error {
	record := si.event["Records"].([]interface{})[0]
	image, ok := record.(map[string]interface{})["dynamodb"].(map[string]interface{})["NewImage"].(map[string]interface{})
	if !ok {
		logger.WithFields(logger.Fields{
			"record": record,
		}).Error("StreamInput::ParseNewImage() missing NewImage attribute in event")
		return errors.New("missing NewImage attribute in event")
	}

	if err := si.unmarshalAttributes(image, out); err != nil {
		return err
	}
	return nil
}

func (si *StreamInput) unmarshalAttributes(attributes map[string]interface{}, out interface{}) error {
	encoded, err := json.Marshal(si.recursivelyFlattenStreamAttributes(attributes))
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err,
		}).Error("StreamInput::unmarshalAttributes() could not marshal json")
		return errors.New("could not marshal json")
	}

	if err := json.Unmarshal(encoded, out); err != nil {
		logger.WithFields(logger.Fields{
			"error":   err,
			"encoded": string(encoded),
		}).Error("StreamInput::unmarshalAttributes() could not unmarshal json")
		return errors.New("could not unmarshal json")
	}
	return nil
}

// Recursively flattens DynamoDB stream attributes into something Go can marshal/unmarshal
func (si *StreamInput) recursivelyFlattenStreamAttributes(attributes map[string]interface{}) map[string]interface{} {
	tmp := make(map[string]interface{})
	for val := range si.flattenStreamAttributes(attributes) {
		tmp[val[0].(string)] = val[1]
		if v, ok := val[1].(map[string]interface{}); ok {
			tmp[val[0].(string)] = si.recursivelyFlattenStreamAttributes(v)
		}
	}
	return tmp
}

// Flattens DynamoDB stream image attributes into something Go can marshal/unmrshal
func (si *StreamInput) flattenStreamAttributes(attributes map[string]interface{}) <-chan []interface{} {
	ch := make(chan []interface{})
	go func() {
		for key, value := range attributes {
			for k, v := range value.(map[string]interface{}) {
				// Stream input format ints as strings so we need to cast them back to ints
				if k == "N" {
					v, _ = strconv.Atoi(v.(string))
				}
				ch <- []interface{}{key, v}
			}
		}
		close(ch)
	}()
	return ch
}

// StreamEventType type of stream event. Possible values: INSERT, MODIFY, REMOVE
type StreamEventType string

// EventType of current stream event
func (si *StreamInput) EventType() StreamEventType {
	record := si.event["Records"].([]interface{})[0]
	return StreamEventType(record.(map[string]interface{})["eventName"].(string))
}

// S3Input ...
type S3Input struct {
	event map[string]interface{}
}

// ObjectKeyPath extract full object key path
func (si *S3Input) ObjectKeyPath() string {
	record := si.event["Records"].([]interface{})[0].(map[string]interface{})
	return record["s3"].(map[string]interface{})["object"].(map[string]interface{})["key"].(string)
}

// ObjectKey extract object key from object key path
func (si *S3Input) ObjectKey() string {
	fragments := strings.Split(si.ObjectKeyPath(), "/")
	return fragments[len(fragments)-1]
}
