package api

import (
    "errors"
    "strings"
    "fmt"
    "log"
)

type Request struct {
    resource      string
    method        string
    event         map[string]interface{}
    handlerConfig *HandlerConfig
}

type HandlerConfig struct {
    Http      map[string]map[string]func(i *Input) *Response
    Scheduled map[string]func() *Response
}

func NewRequest(event interface{}, config *HandlerConfig) *Request {
    resource, httpMethod := "", ""
    if val, ok := event.(map[string]interface{})["resource"]; ok {
        resource = val.(string)
    }

    if val, ok := event.(map[string]interface{})["httpMethod"]; ok {
        httpMethod = val.(string)
    }

    return &Request{
        resource:      resource,
        method:        httpMethod,
        event:         event.(map[string]interface{}),
        handlerConfig: config,
    }
}

// Invoke invoke and handle request by event type.
// Supported events are: Api and Schedule
func (r *Request) Invoke() (*Response, error) {
    log.Printf("Request event: %s", r.event)

    switch true {
    case r.isHttpEvent():
        return r.handleHttpEvent()
    case r.isScheduledEvent():
        return r.handleScheduledEvent()
    }
    return nil, errors.New("unknown event")
}

func (r Request) isHttpEvent() bool {
    if _, ok := r.event["path"]; ok {
        return true
    }
    return false
}

func (r Request) isScheduledEvent() bool {
    return r.event["type"] == "schedule"
}

func (r *Request) handleHttpEvent() (*Response, error) {
    pathParams, ok := r.event["pathParameters"]
    resource := r.resource

    if ok && pathParams != nil {
        for k, v := range pathParams.(map[string]interface{}) {
            resource = strings.Replace(resource, v.(string), fmt.Sprintf("{%s}", k), 1)
        }
    }

    handler, found := r.handlerConfig.Http[resource][r.method]
    if !found {
        return nil, errors.New("handler func missing")
    }

    return &*handler(&Input{event: r.event}), nil
}

func (r *Request) handleScheduledEvent() (*Response, error) {
    handler, found := r.handlerConfig.Scheduled[r.resource]
    if !found {
        return nil, errors.New("handler func missing")
    }

    return &*handler(), nil
}
