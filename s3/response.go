package s3

import "github.com/matthisstenius/logger"

// Response for S3 event
type Response struct{}

// NewResponse initializer
func NewResposne(message string) *Response {
	logger.WithFields(logger.Fields{"message": message}).Info("Response::NewResponse() schedule handler responded")
	return &Response{}
}
