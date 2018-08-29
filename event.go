package router

import (
	"errors"
	"runtime/debug"

	"github.com/matthisstenius/lambda-router/domain"
	"github.com/matthisstenius/logger"
)

// Event ...
type Event struct {
	config *Config
}

// Config for routing event handlers
type Config struct {
	HTTP      domain.Router
	Scheduled domain.Router
	DynamoDB  domain.Router
	S3        domain.Router
	SNS       domain.Router
}

// NewEvent initialization for Event
func NewEvent(config *Config) *Event {
	return &Event{config: config}
}

// Handle event by routing matched event
func (e *Event) Handle(event interface{}) (interface{}, error) {
	evt := event.(map[string]interface{})
	logger.WithFields(logger.Fields{
		"event": event,
	}).Info("Incoming event")
	defer e.logPanic()

	var response domain.Response
	var err error
	switch true {
	case e.config.HTTP.IsMatch(evt):
		response, err = e.config.HTTP.Route(evt)
		break
	case e.config.Scheduled.IsMatch(evt):
		response, err = e.config.Scheduled.Route(evt)
		break
	case e.config.DynamoDB.IsMatch(evt):
		response, err = e.config.DynamoDB.Route(evt)
		break
	case e.config.S3.IsMatch(evt):
		response, err = e.config.S3.Route(evt)
		break
	case e.config.SNS.IsMatch(evt):
		response, err = e.config.SNS.Route(evt)
		break
	default:
		response, err = nil, errors.New("unknown event")
	}

	if response != nil {
		return response.Payload(), err
	}
	return nil, err
}

func (e *Event) logPanic() {
	if r := recover(); r != nil {
		logger.WithFields(logger.Fields{
			"error": r,
			"stack": string(debug.Stack()),
		}).Error("Unexpected panic")
	}
}
