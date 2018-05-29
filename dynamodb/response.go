package dynamodb

import "github.com/matthisstenius/logger"

// Response for dynamodb event
type Response struct{}

// NewResponse initializer
func NewResponse(message string) *Response {
	logger.WithFields(logger.Fields{"message": message}).Info("Response::NewResponse() dynamodb handler responded")
	return &Response{}
}
