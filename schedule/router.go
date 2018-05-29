package schedule

import "errors"

// Routes mappings for Schedule handlers
type Routes map[string]func() *Response

// Router for Schedule events
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
	resource := r.event["resource"].(string)
	handler, found := r.routes[resource]
	if !found {
		return nil, errors.New("handler func missing")
	}

	return &*handler(), nil
}

// IsMatch for Schedule event
func IsMatch(e map[string]interface{}) bool {
	return e["type"] == "schedule"
}
