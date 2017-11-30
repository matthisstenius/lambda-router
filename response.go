package api

import (
	"encoding/json"
)

type Response struct {
	StatusCode int
	Body interface{}
}

func NewResponse(status int, body interface{}) Response {
	encoded, _ := json.Marshal(map[string]interface{}{"data": body})
	return Response{StatusCode: status, Body: string(encoded)}
}
