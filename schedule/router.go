package schedule

import (
	"errors"
	"github.com/matthisstenius/lambda-router/v4/domain"
)

const EventSource = "schedule"

// Route mapping for handler
type Route struct {
	Handler func() domain.Response
}

// Routes mappings for Schedule handlers
type Routes map[string]Route

// Router for Schedule events
type Router struct {
	routes Routes
}

// NewRouter initializer
func NewRouter(routes Routes) *Router {
	return &Router{routes: routes}
}

// Route incoming event to corresponding handler
func (r *Router) Route(evt map[string]interface{}) (domain.Response, error) {
	resource := evt["resource"].(string)
	route, found := r.routes[resource]
	if !found {
		return nil, errors.New("handler func missing")
	}
	return route.Handler(), nil
}

// IsMatch for Schedule event
func (r *Router) IsMatch(e map[string]interface{}) bool {
	return e["eventSource"] == EventSource
}
