package http

import (
	"errors"
	"fmt"
	"strings"
)

// Routes mappings for HTTP handlers
type Routes map[string]map[string]func(i *Input) *Response

// Router for HTTP events
type Router struct {
	event  map[string]interface{}
	routes Routes
}

// NewRouter initializer
func NewRouter(e map[string]interface{}, routes Routes) *Router {
	return &Router{event: e, routes: routes}
}

// Dispatch incoming event to corresponding handler
func (r *Router) Dispatch() (*Response, error) {
	pathParams, ok := r.event["pathParameters"]
	resource := r.event["resource"].(string)
	method := r.event["httpMethod"].(string)

	if ok && pathParams != nil {
		for k, v := range pathParams.(map[string]interface{}) {
			resource = strings.Replace(resource, v.(string), fmt.Sprintf("{%s}", k), 1)
		}
	}

	handler, ok := r.routes[resource][method]
	if !ok {
		return nil, errors.New("handler func missing")
	}
	return &*handler(NewInput(r.event)), nil
}

// IsMatch for HTTP event
func IsMatch(e map[string]interface{}) bool {
	if _, ok := e["httpMethod"]; ok {
		return true
	}
	return false
}
