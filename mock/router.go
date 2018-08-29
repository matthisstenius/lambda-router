package mock

import "github.com/matthisstenius/lambda-router/domain"

// Router mock
type Router struct {
	DispatchFn func(evt map[string]interface{}) (domain.Response, error)
	IsMatchFn  func(evt map[string]interface{}) bool
}

// Route mock implementation
func (r *Router) Route(evt map[string]interface{}) (domain.Response, error) {
	return r.DispatchFn(evt)
}

// IsMatch mock implementation
func (r *Router) IsMatch(evt map[string]interface{}) bool {
	return r.IsMatchFn(evt)
}
