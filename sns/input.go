package sns

import (
	"encoding/json"
	"errors"

	"github.com/matthisstenius/logger"
)

// Input for parsed SNS event
// TODO: Write tests
type Input struct {
	event map[string]interface{}
}

// NewInput initializer
func NewInput(e map[string]interface{}) *Input {
	return &Input{event: e}
}

// ParseMessage as JSON
func (i *Input) ParseMessage(out interface{}) error {
	record := i.event["Records"].([]interface{})[0].(map[string]interface{})
	if err := json.Unmarshal([]byte(record["Sns"].(map[string]interface{})["Message"].(string)), out); err != nil {
		logger.WithFields(logger.Fields{
			"error": err,
		}).Error("SNSInput::ParseMessage() could not unmarshal json")
		return errors.New("invalid SNS payload")
	}
	return nil
}
