package cognito

import (
	"errors"
)

// AccessProvider cognito implementation
type AccessProvider struct{}

// ParseRoles parse roles from given event
func (ap *AccessProvider) ParseRoles(evt map[string]interface{}) ([]string, error) {
	reqContext := evt["requestContext"]
	authorizer, ok := reqContext.(map[string]interface{})["authorizer"]
	if !ok || authorizer == nil {
		return nil, errors.New("authorizer index missing in event")
	}

	claims, ok := authorizer.(map[string]interface{})["claims"]
	if !ok {
		return nil, errors.New("claims index missing in authorizer")
	}
	roles, ok := claims.(map[string]interface{})["cognito:groups"]
	if !ok {
		return nil, errors.New("cognito:groups index missing in claims")
	}
	return []string{roles.(string)}, nil
}
