package http

import (
	"encoding/json"

	"github.com/matthisstenius/lambda-router/domain"
	"github.com/matthisstenius/logger"
)

// Response for HTTP event
type Response struct {
	StatusCode      int               `json:"statusCode"`
	Body            interface{}       `json:"body"`
	Headers         map[string]string `json:"headers"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) Success(status int, body interface{}) domain.Response {
	encoded, _ := json.Marshal(body)
	logger.WithFields(logger.Fields{"body": string(encoded)}).Info("response")

	return &Response{
		StatusCode: status,
		Body:       string(encoded),
		Headers: map[string]string{
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Origin":  "*",
		},
		IsBase64Encoded: false,
	}
}

func (r *Response) Error(status int, error interface{}) domain.Response {
	encoded, _ := json.Marshal(map[string]interface{}{
		"error": error,
	})
	logger.WithFields(logger.Fields{
		"error": string(encoded),
	}).Info("Error response")

	return &Response{
		StatusCode: status,
		Body:       string(encoded),
		Headers: map[string]string{
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Origin":  "*",
		},
		IsBase64Encoded: false,
	}
}

//// NewResponse initialize success response
//func NewResponse(status int, body interface{}) *Response {
//	encoded, _ := json.Marshal(body)
//	logger.WithFields(logger.Fields{"body": string(encoded)}).Info("response")
//
//	return &Response{
//		StatusCode: status,
//		Body:       string(encoded),
//		Headers: map[string]string{
//			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key",
//			"Access-Control-Allow-Methods": "*",
//			"Access-Control-Allow-Origin":  "*",
//		},
//		IsBase64Encoded: false,
//	}
//}
//
//// NewErrorResponse initialize error response
//func NewErrorResponse(status int, error interface{}) *Response {
//	encoded, _ := json.Marshal(map[string]interface{}{
//		"error": error,
//	})
//	logger.WithFields(logger.Fields{
//		"error": string(encoded),
//	}).Info("Error response")
//
//	return &Response{
//		StatusCode: status,
//		Body:       string(encoded),
//		Headers: map[string]string{
//			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key",
//			"Access-Control-Allow-Methods": "*",
//			"Access-Control-Allow-Origin":  "*",
//		},
//		IsBase64Encoded: false,
//	}
//}
