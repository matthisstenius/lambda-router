package sns

import (
	"errors"
	"github.com/matthisstenius/lambda-router/v4/domain"
)

const EventSource = "aws:sns"

// Route mapping for handler
type Route struct {
	Handler func(i *Input) domain.Response
}

// Routes mappings for SNS handlers
type Routes map[string]Route

// Router for SNS events
type Router struct {
	routes Routes
}

// NewRouter initializer
func NewRouter(routes Routes) *Router {
	return &Router{routes: routes}
}

// Route incoming event to corresponding handler
func (r *Router) Route(evt map[string]interface{}) (domain.Response, error) {
	record := evt["Records"].([]interface{})[0].(map[string]interface{})
	route, ok := r.routes[record["Sns"].(map[string]interface{})["TopicArn"].(string)]
	if !ok {
		return nil, errors.New("handler func missing")
	}
	return route.Handler(NewInput(evt)), nil
}

// IsMatch for SNS event
func (r *Router) IsMatch(e map[string]interface{}) bool {
	if v, ok := e["Records"].([]interface{}); ok && len(v) > 0 {
		return v[0].(map[string]interface{})["EventSource"] == EventSource
	}
	return false
}
