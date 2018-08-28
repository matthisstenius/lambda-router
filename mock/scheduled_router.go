package mock

import "github.com/matthisstenius/lambda-router/domain"

// ScheduledRouter mock
type ScheduledRouter struct {
	DispatchFn func(evt map[string]interface{}) (domain.Response, error)
	IsMatchFn  func(evt map[string]interface{}) bool
}

// Route mock implementation
func (sr *ScheduledRouter) Route(evt map[string]interface{}) (domain.Response, error) {
	return sr.DispatchFn(evt)
}

// IsMatch mock implementation
func (sr *ScheduledRouter) IsMatch(evt map[string]interface{}) bool {
	return sr.IsMatchFn(evt)
}
