package api

import (
	"errors"
)

type Request struct {
	resource string
	method string
	routes map[string]map[string]func()Response
}

func NewRequest(event interface{}, routes map[string]map[string]func()Response) Request {
	return Request{
		resource: event.(map[string]interface{})["resource"].(string),
		method: event.(map[string]interface{})["httpMethod"].(string),
		routes: routes,
	}
}

func (r Request) Invoke() (Response, error) {
	handler, ok := r.routes[r.resource][r.method]

	var response Response
	if !ok {
		return response, errors.New("handler func missing")
	}

	response = handler()
	return response, nil
}