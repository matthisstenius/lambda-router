package lambdaRouter

import (
	"errors"
	"encoding/json"
)

type Router struct {
	routes map[string]map[string]interface{}
	request Request
}
func New(routes map[string]map[string]interface{}, request Request) *Router {
	return &Router{routes: routes, request: request}
}

func (r Router) Invoke() (string, error) {
	handler, ok := r.routes[r.request.resource][r.request.method]

	if !ok {
		return "", errors.New("handler func missing")
	}

	data := handler.(func() interface{})()
	encoded, err := json.Marshal(data)

	if err != nil {
		return "", errors.New("invalid response data")
	}

	return string(encoded), nil
}
