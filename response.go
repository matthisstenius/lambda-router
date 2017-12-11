package api

import (
	"encoding/json"
)

type Response struct {
	StatusCode      int               `json:"statusCode"`
	Body            interface{}       `json:"body"`
	Headers         map[string]string `json:"headers"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

func NewResponse(status int, body interface{}) Response {
	encoded, _ := json.Marshal(map[string]interface{}{"data": body})
	return Response{
		StatusCode:      status,
		Body:            string(encoded),
		Headers:         map[string]string{},
		IsBase64Encoded: false,
	}
}
