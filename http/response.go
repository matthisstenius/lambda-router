package http

import (
	"encoding/json"
	"github.com/matthisstenius/logger"
)

// Response for HTTP event
type Response struct {
	statusCode      int
	body            interface{}
	headers         map[string]string
	isBase64Encoded bool
}

func (r *Response) Payload() interface{} {
	return map[string]interface{}{
		"statusCode":      r.statusCode,
		"body":            r.body,
		"headers":         r.headers,
		"isBase64Encoded": r.isBase64Encoded,
	}
}

// NewResponse initialize success response
func NewResponse(status int, body interface{}) *Response {
	encoded, _ := json.Marshal(body)
	logger.WithFields(logger.Fields{"body": string(encoded)}).Info("response")

	return &Response{
		statusCode: status,
		body:       string(encoded),
		headers: map[string]string{
			"Access-Control-Allow-headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Origin":  "*",
		},
		isBase64Encoded: false,
	}
}

// NewErrorResponse initialize error response
func NewErrorResponse(status int, error interface{}) *Response {
	encoded, _ := json.Marshal(map[string]interface{}{
		"error": error,
	})
	logger.WithFields(logger.Fields{
		"error": string(encoded),
	}).Info("Error response")

	return &Response{
		statusCode: status,
		body:       string(encoded),
		headers: map[string]string{
			"Access-Control-Allow-headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Origin":  "*",
		},
		isBase64Encoded: false,
	}
}
