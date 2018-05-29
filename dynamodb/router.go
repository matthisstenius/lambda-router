package dynamodb

import (
	"errors"
)

const eventSource = "aws:dynamodb"

// Routes mappings for DynamoDB handlers
type Routes map[string]func(i *Input) *Response

// Router for DynamoDB events
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
	record := r.event["Records"].([]interface{})[0]
	streamArn := record.(map[string]interface{})["eventSourceARN"].(string)
	handler, ok := r.routes[streamArn]
	if ok {
		return nil, errors.New("handler func missing")
	}
	return &*handler(&Input{event: r.event}), nil
}

// IsMatch for DynamoDB event
func IsMatch(e map[string]interface{}) bool {
	if v, ok := e["Records"].([]interface{}); ok && len(v) > 0 {
		return v[0].(map[string]interface{})["eventSource"] == eventSource
	}
	return false
}
