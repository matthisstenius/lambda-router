package sns

import (
	"encoding/json"
	"errors"

	"github.com/matthisstenius/logger"
)

// Input for parsed SNS event
type Input struct {
	event map[string]interface{}
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
