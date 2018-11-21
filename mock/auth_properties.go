package mock

// AuthProperties mock implementation
type AuthProperties struct{}

// GetParam mock implementation
func (ap *AuthProperties) GetParam(key string) interface{} {
	return nil
}

// HasRole mock implementation
func (ap *AuthProperties) HasRole(role string) bool {
	return false
}
