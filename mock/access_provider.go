package mock

// AccessProvider mock
type AccessProvider struct {
	RolesFn func(evt map[string]interface{}) ([]string, error)
}

// Roles mock implementation
func (ap *AccessProvider) Roles(evt map[string]interface{}) ([]string, error) {
	return ap.RolesFn(evt)
}
