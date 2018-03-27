package api

import (
    "errors"
    "strings"
    "fmt"
    "log"
    "github.com/sirupsen/logrus"
    "runtime/debug"
)

type Handler struct {
    config *HandlerConfig
    event  map[string]interface{}
}

type HandlerConfig struct {
    Http      map[string]map[string]func(i *Input) *Response
    Scheduled map[string]func() *Response
}

func Newhandler(config *HandlerConfig) *Handler {
    return &Handler{config: config}
}

func (h *Handler) Invoke(event interface{}) (*Response, error) {
    h.event = event.(map[string]interface{})
    log.Printf("Request event: %s", event)
    defer h.logPanic()

    switch true {
    case h.isHttpEvent():
        return h.handleHttpEvent()
    case h.isScheduledEvent():
        return h.handleScheduledEvent()
    }
    return nil, errors.New("unknown event")
}

func (h *Handler) isHttpEvent() bool {
    if _, ok := h.event["httpMethod"]; ok {
        return true
    }
    return false
}

func (h *Handler) isScheduledEvent() bool {
    return h.event["type"] == "schedule"
}

func (h *Handler) handleHttpEvent() (*Response, error) {
    pathParams, ok := h.event["pathParameters"]
    resource := h.event["resource"].(string)
    method := h.event["httpMethod"].(string)

    if ok && pathParams != nil {
        for k, v := range pathParams.(map[string]interface{}) {
            resource = strings.Replace(resource, v.(string), fmt.Sprintf("{%s}", k), 1)
        }
    }

    handler, found := h.config.Http[resource][method]
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

func (h *Handler) logPanic() {
    if r := recover(); r != nil {
        logrus.SetFormatter(&logrus.JSONFormatter{})
        logrus.WithFields(logrus.Fields{
            "error": r,
            "stack": string(debug.Stack()),
        }).Error()
    }
}
