package api

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"regexp"

	"encoding/json"

	"bitbucket.org/mstenius/logger"
)

const (
	eventSourceDynamoDB = "aws:dynamodb"
	eventSourceS3       = "aws:s3"
	eventSourceSNS      = "aws:sns"
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
	S3        map[string]func(i *S3Input) *Response
	SNS       map[string]func(i *SNSInput) *Response
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
		break
	case h.isS3Event():
		response, err = h.handleS3Event()
		break
	case h.isSNSEvent():
		response, err = h.handleSNSEvent()
		break
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
	if v, ok := h.event["Records"].([]interface{}); ok && len(v) > 0 {
		return v[0].(map[string]interface{})["eventSource"] == eventSourceDynamoDB
	}
	return false
}

func (h *Handler) isS3Event() bool {
	if v, ok := h.event["Records"].([]interface{}); ok && len(v) > 0 {
		return v[0].(map[string]interface{})["eventSource"] == eventSourceS3
	}
	return false
}

func (h *Handler) isSNSEvent() bool {
	if v, ok := h.event["Records"].([]interface{}); ok && len(v) > 0 {
		return v[0].(map[string]interface{})["eventSource"] == eventSourceSNS
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
	record := h.event["Records"].([]interface{})[0]
	streamArn := record.(map[string]interface{})["eventSourceARN"].(string)
	handler, ok := h.config.Stream[streamArn]
	if !ok {
		return nil, errors.New("handler func missing")
	}
	return &*handler(&StreamInput{event: h.event}), nil
}

func (h *Handler) handleS3Event() (*Response, error) {
	record := h.event["Records"].([]interface{})[0].(map[string]interface{})
	key := record["s3"].(map[string]interface{})["object"].(map[string]interface{})["key"].(string)

	re := regexp.MustCompile("[^/]+$")
	folder := re.ReplaceAllString(key, "")
	if folder == "" {
		// If object is in root we want to look to /
		folder = "/"
	} else {
		folder = strings.TrimSuffix(folder, "/")
	}

	handler, ok := h.config.S3[folder]
	if !ok {
		return nil, errors.New("handler func missing")
	}
	return &*handler(&S3Input{event: h.event}), nil
}

func (h *Handler) handleSNSEvent() (*Response, error) {
	record := h.event["Records"].([]interface{})[0].(map[string]interface{})
	var message map[string]interface{}
	if err := json.Unmarshal([]byte(record["Sns"].(string)), &message); err != nil {
		return nil, errors.New("invalid SNS payload")
	}

	handler, ok := h.config.SNS[message["messageType"].(string)]
	if !ok {
		return nil, errors.New("handler func missing")
	}
	return &*handler(&SNSInput{event: h.event}), nil
}

func (h *Handler) logPanic() {
	if r := recover(); r != nil {
		logger.WithFields(logger.Fields{
			"error": r,
			"stack": string(debug.Stack()),
		}).Error("Unexpected panic")
	}
}
