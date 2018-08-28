package mock

import "github.com/matthisstenius/lambda-router/domain"

// HTTPRouter mock
type HTTPRouter struct {
	DispatchFn func(evt map[string]interface{}) (domain.Response, error)
	IsMatchFn  func(evt map[string]interface{}) bool
}

// Route mock implementation
func (hr *HTTPRouter) Route(evt map[string]interface{}) (domain.Response, error) {
	return hr.DispatchFn(evt)
}

// IsMatch mock implementation
func (hr *HTTPRouter) IsMatch(evt map[string]interface{}) bool {
	return hr.IsMatchFn(evt)
}
