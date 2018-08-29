package dynamodb

import (
	"errors"
	"github.com/matthisstenius/lambda-router/v2/domain"
)

const EventSource = "aws:dynamodb"

// Route mapping for handler and optional access
type Route struct {
	Handler func(i *Input) domain.Response
}

// Routes mappings for HTTP handlers
type Routes map[string]Route

// Router for DynamoDB events
type Router struct {
	routes Routes
}

// NewRouter initializer
func NewRouter(routes Routes) *Router {
	return &Router{routes: routes}
}

// Dispatch incoming event to corresponding handler
func (r *Router) Route(evt map[string]interface{}) (domain.Response, error) {
	record := evt["Records"].([]interface{})[0]
	streamArn := record.(map[string]interface{})["eventSourceARN"].(string)

	route, ok := r.routes[streamArn]
	if !ok {
		return nil, errors.New("handler func missing")
	}
	return route.Handler(NewInput(evt)), nil
}

// IsMatch for DynamoDB event
func (r *Router) IsMatch(e map[string]interface{}) bool {
	if v, ok := e["Records"].([]interface{}); ok && len(v) > 0 {
		return v[0].(map[string]interface{})["eventSource"] == EventSource
	}
	return false
}
