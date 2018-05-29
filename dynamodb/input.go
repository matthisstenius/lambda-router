package dynamodb

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/matthisstenius/logger"
)

// Input for parsed DynamoDB event
type Input struct {
	event map[string]interface{}
}

// ParseOldImage from DynamoDB dynamodb event
func (i *Input) ParseOldImage(out interface{}) error {
	record := i.event["Records"].([]interface{})[0]
	image, ok := record.(map[string]interface{})["dynamodb"].(map[string]interface{})["OldImage"].(map[string]interface{})
	if !ok {
		logger.WithFields(logger.Fields{
			"record": record,
		}).Error("StreamInput::ParseOldImage() missing OldImage attribute in event")
		return errors.New("missing OldImage attribute in event")
	}

	if err := i.unmarshalAttributes(image, out); err != nil {
		return err
	}
	return nil
}

// ParseNewImage from DynamoDB dynamodb event
func (i *Input) ParseNewImage(out interface{}) error {
	record := i.event["Records"].([]interface{})[0]
	image, ok := record.(map[string]interface{})["dynamodb"].(map[string]interface{})["NewImage"].(map[string]interface{})
	if !ok {
		logger.WithFields(logger.Fields{
			"record": record,
		}).Error("StreamInput::ParseNewImage() missing NewImage attribute in event")
		return errors.New("missing NewImage attribute in event")
	}

	if err := i.unmarshalAttributes(image, out); err != nil {
		return err
	}
	return nil
}

func (i *Input) unmarshalAttributes(attributes map[string]interface{}, out interface{}) error {
	encoded, err := json.Marshal(i.recursivelyFlattenStreamAttributes(attributes))
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

// Recursively flattens DynamoDB dynamodb attributes into something Go can marshal/unmarshal
func (i *Input) recursivelyFlattenStreamAttributes(attributes map[string]interface{}) map[string]interface{} {
	tmp := make(map[string]interface{})
	for val := range i.flattenStreamAttributes(attributes) {
		tmp[val[0].(string)] = val[1]
		if v, ok := val[1].(map[string]interface{}); ok {
			tmp[val[0].(string)] = i.recursivelyFlattenStreamAttributes(v)
		}
	}
	return tmp
}

// Flattens DynamoDB dynamodb image attributes into something Go can marshal/unmrshal
func (i *Input) flattenStreamAttributes(attributes map[string]interface{}) <-chan []interface{} {
	ch := make(chan []interface{})
	go func() {
		for key, value := range attributes {
			for k, v := range value.(map[string]interface{}) {
				// Stream dynamodb format ints as strings so we need to cast them back to ints
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

const (
	EventInsert EventType = "INSERT"
	EventModify EventType = "MODIFY"
	EventRemove EventType = "REMOVE"
)

// EventType type of dynamodb event. Possible values: INSERT, MODIFY, REMOVE
type EventType string

// EventType of current dynamodb event
func (i *Input) EventType() EventType {
	record := i.event["Records"].([]interface{})[0]
	return EventType(record.(map[string]interface{})["eventName"].(string))
}
