package mock

import "github.com/matthisstenius/lambda-router/v2/domain"

// AuthProvider mock implementation
type AuthProvider struct{}

// ParseAuth mock implementation
func (ap *AuthProvider) ParseAuth(evt map[string]interface{}) (domain.AuthProperties, error) {
	return domain.AuthProperties{}, nil
}
