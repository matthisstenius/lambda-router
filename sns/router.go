package sns

import "errors"

const eventSourceSNS = "aws:sns"

// Routes mappings for SNS handlers
type Routes map[string]func(i *Input) *Response

// Router for SNS events
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
	record := r.event["Records"].([]interface{})[0].(map[string]interface{})
	handler, ok := r.routes[record["Sns"].(map[string]interface{})["TopicArn"].(string)]
	if !ok {
		return nil, errors.New("handler func missing")
	}
	return &*handler(NewInput(r.event)), nil
}

// IsMatch for SNS event
func IsMatch(e map[string]interface{}) bool {
	if v, ok := e["Records"].([]interface{}); ok && len(v) > 0 {
		return v[0].(map[string]interface{})["EventSource"] == eventSourceSNS
	}
	return false
}
