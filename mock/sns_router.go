package mock

import "github.com/matthisstenius/lambda-router/domain"

// SNSRouter mock
type SNSRouter struct {
	DispatchFn func(evt map[string]interface{}) (domain.Response, error)
	IsMatchFn  func(evt map[string]interface{}) bool
}

// Route mock implementation
func (sr *SNSRouter) Route(evt map[string]interface{}) (domain.Response, error) {
	return sr.DispatchFn(evt)
}

// IsMatch mock implementation
func (sr *SNSRouter) IsMatch(evt map[string]interface{}) bool {
	return sr.IsMatchFn(evt)
}
