package api

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"bitbucket.org/mstenius/logger"
)

const (
	eventSourceDynamoDB = "aws:dynamodb"
)

// Handler ...
type Handler struct {
	config *HandlerConfig
	event  map[string]interface{}
}

// HandlerConfig for event handlers
type HandlerConfig struct {
	HTTP      map[string]map[string]func(i *Input) *Response
	Scheduled map[string]func() *Response
	Stream    map[string]func(i *StreamInput) *Response
}

// NewHandler initialization for Handler
func NewHandler(config *HandlerConfig) *Handler {
	return &Handler{config: config}
}

// Invoke correct handler based in mapped handler config and incoming event
func (h *Handler) Invoke(event interface{}) (*Response, error) {
	h.event = event.(map[string]interface{})
	logger.WithFields(logger.Fields{
		"event": event,
	}).Info("Incoming event")
	defer h.logPanic()

	var response *Response
	var err error
	switch true {
	case h.isHTTPEvent():
		response, err = h.handleHTTPEvent()
		break
	case h.isScheduledEvent():
		response, err = h.handleScheduledEvent()
		break
	case h.isStreamEvent():
		response, err = h.handleStreamEvent()
	default:
		response, err = nil, errors.New("unknown event")
	}
	return response, err
}

func (h *Handler) isHTTPEvent() bool {
	if _, ok := h.event["httpMethod"]; ok {
		return true
	}
	return false
}

func (h *Handler) isScheduledEvent() bool {
	return h.event["type"] == "schedule"
}

func (h *Handler) isStreamEvent() bool {
	if v, ok := h.event["Records"]; ok {
		if len(v.([]map[string]interface{})) > 0 && v.([]map[string]string)[0]["EventSource"] == eventSourceDynamoDB {
			return true
		}
	}
	return false
}

func (h *Handler) handleHTTPEvent() (*Response, error) {
	pathParams, ok := h.event["pathParameters"]
	resource := h.event["resource"].(string)
	method := h.event["httpMethod"].(string)

	if ok && pathParams != nil {
		for k, v := range pathParams.(map[string]interface{}) {
			resource = strings.Replace(resource, v.(string), fmt.Sprintf("{%s}", k), 1)
		}
	}

	handler, found := h.config.HTTP[resource][method]
	if !found {
		return nil, errors.New("handler func missing")
	}

	return &*handler(&Input{event: h.event}), nil
}

func (h *Handler) handleScheduledEvent() (*Response, error) {
	resource := h.event["resource"].(string)
	handler, found := h.config.Scheduled[resource]
	if !found {
		return nil, errors.New("handler func missing")
	}

	return &*handler(), nil
}

func (h *Handler) handleStreamEvent() (*Response, error) {
	record := h.event["Records"].([]map[string]interface{})[0]
	streamArn := record["eventSourceARN"].(string)
	handler, ok := h.config.Stream[streamArn]
	if !ok {
		return nil, errors.New("handler func missing")
	}
	return &*handler(&StreamInput{event: h.event}), nil
}

func (h *Handler) logPanic() {
	if r := recover(); r != nil {
		logger.WithFields(logger.Fields{
			"error": r,
			"stack": string(debug.Stack()),
		}).Error("Unexpected panic")
	}
}
