package api

import (
	"errors"
	"strings"
	"fmt"
	"log"
)

type Request struct {
	resource string
	method   string
	event    map[string]interface{}
	routes   map[string]map[string]func(i *Input) *Response
}

func NewRequest(event interface{}, routes map[string]map[string]func(i *Input) *Response) *Request {
	return &Request{
		resource: event.(map[string]interface{})["resource"].(string),
		method:   event.(map[string]interface{})["httpMethod"].(string),
		event:    event.(map[string]interface{}),
		routes:   routes,
	}
}

func (r *Request) Invoke() (*Response, error) {
	log.Printf("Request event: %s", r.event)
	resource := r.resource
	pathParams, ok := r.event["pathParameters"]

	if ok && pathParams != nil {
		for k, v := range pathParams.(map[string]interface{}) {
			resource = strings.Replace(resource, v.(string), fmt.Sprintf("{%s}", k), 1)
		}
	}

	handler, ok := r.routes[resource][r.method]

	var response Response
	if !ok {
		return &response, errors.New("handler func missing")
	}

	response = *handler(&Input{event: r.event})
	return &response, nil
}
