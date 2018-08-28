package mock

import "github.com/matthisstenius/lambda-router/domain"

// DynamoRouter mock
type DynamoRouter struct {
	DispatchFn func(evt map[string]interface{}) (domain.Response, error)
	IsMatchFn  func(evt map[string]interface{}) bool
}

// Route mock implementation
func (dr *DynamoRouter) Route(evt map[string]interface{}) (domain.Response, error) {
	return dr.DispatchFn(evt)
}

// IsMatch mock implementation
func (dr *DynamoRouter) IsMatch(evt map[string]interface{}) bool {
	return dr.IsMatchFn(evt)
}
