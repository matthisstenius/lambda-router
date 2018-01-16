package api

import (
    "errors"
    "strings"
    "fmt"
    "log"
    "encoding/json"
)

type Request struct {
    resource string
    method   string
    event    map[string]interface{}
    routes   map[string]map[string]func(i *Input) *Response
    events map[string]func(i *Input) *Response
}

func NewRequest(event interface{}, routes map[string]map[string]func(i *Input) *Response, events map[string]func(i *Input) *Response) *Request {
    return &Request{
        resource: event.(map[string]interface{})["resource"].(string),
        method:   event.(map[string]interface{})["httpMethod"].(string),
        event:    event.(map[string]interface{}),
        routes:   routes,
        events: events,
    }
}

func (r *Request) Invoke() (*Response, error) {
    log.Printf("Request event: %s", r.event)
    var handler func(i *Input) *Response
    var found bool

    if !r.isSnsEvent() {
        resource := r.resource
        pathParams, ok := r.event["pathParameters"]
        if ok && pathParams != nil {
            for k, v := range pathParams.(map[string]interface{}) {
                resource = strings.Replace(resource, v.(string), fmt.Sprintf("{%s}", k), 1)
            }
        }
        handler, found = r.routes[resource][r.method]
    } else {
        record := r.event["Records"].([]map[string]interface{})[0]
        message := record["SNS"].(map[string]string)["Message"]

        var data map[string]interface{}
        if err := json.Unmarshal([]byte(message), &data); err != nil {
            return nil, errors.New("could not parse SNS Message")
        }

        handler, found = r.events[data["messageType"].(string)]
    }

    var response Response
    if !found {
        return &response, errors.New("handler func missing")
    }

    response = *handler(&Input{event: r.event})
    return &response, nil
}

func (r Request) isSnsEvent() bool {
    return len(r.event["Records"].([]map[string]interface{})) > 0
}