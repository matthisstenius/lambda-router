package lambdaRouter

import (
	"errors"
)

type Router struct {
	routes map[string]map[string]interface{}
	request Request
}
func New(routes map[string]map[string]interface{}, event interface{}) *Router {
	request := Request{
		resource: event.(map[string]interface{})["resource"].(string),
		method: event.(map[string]interface{})["httpMethod"].(string),
	}
	return &Router{routes: routes, request: request}
}

func (r Router) Invoke() (Response, error) {
	handler, ok := r.routes[r.request.resource][r.request.method]

	var response Response
	if !ok {
		return response, errors.New("handler func missing")
	}

	response = handler.(func() Response)()
	return response, nil
}
