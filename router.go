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

func (r Router) Invoke() (interface{}, error) {
	handler, ok := r.routes[r.request.resource][r.request.method]

	if !ok {
		return "", errors.New("handler func missing")
	}

	data := handler.(func() interface{})()

	return data, nil
}
