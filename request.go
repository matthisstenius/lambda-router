package api

import (
	"errors"
)

type Request struct {
	resource string
	method string
	event map[string]interface{}
	routes map[string]map[string]func(i Input)Response
}

func NewRequest(event interface{}, routes map[string]map[string]func(i Input)Response) Request {
	return Request{
		resource: event.(map[string]interface{})["resource"].(string),
		method: event.(map[string]interface{})["httpMethod"].(string),
		event: event.(map[string]interface{}),
		routes: routes,
	}
}

func (r Request) Invoke() (Response, error) {
	handler, ok := r.routes[r.resource][r.method]

	var response Response
	if !ok {
		return response, errors.New("handler func missing")
	}

	response = handler(Input{event: r.event})
	return response, nil
}