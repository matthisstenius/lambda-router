package schedule

import "github.com/matthisstenius/logger"

// Response for schedule event
type Response struct{}

// NewResponse initializer
func NewResposne(message string) *Response {
	logger.WithFields(logger.Fields{"message": message}).Info("Response::NewResponse() schedule handler responded")
	return &Response{}
}
