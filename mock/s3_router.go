package mock

import "github.com/matthisstenius/lambda-router/domain"

// S3Router mock
type S3Router struct {
	DispatchFn func(evt map[string]interface{}) (domain.Response, error)
	IsMatchFn  func(evt map[string]interface{}) bool
}

// Route mock implementation
func (sr *S3Router) Route(evt map[string]interface{}) (domain.Response, error) {
	return sr.DispatchFn(evt)
}

// IsMatch mock implementation
func (sr *S3Router) IsMatch(evt map[string]interface{}) bool {
	return sr.IsMatchFn(evt)
}
