package s3

import "github.com/matthisstenius/logger"

// Response for S3 event
type Response struct{}

// Payload formatted response data
func (r *Response) Payload() interface{} {
	return "No payload"
}

// NewResponse initializer
func NewResponse(message string) *Response {
	logger.WithFields(logger.Fields{"message": message}).Info("Response::NewResponse() S3 handler responded")
	return &Response{}
}
