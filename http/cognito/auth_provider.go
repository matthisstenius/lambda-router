package cognito

import (
	"encoding/json"
	"errors"
	"github.com/matthisstenius/logger"
)

// AuthProvider Cognito implementation
type AuthProvider struct{}

// ParseAuth data from event
func (ap *AuthProvider) ParseAuth(evt map[string]interface{}) (AuthProperties, error) {
	reqContext := evt["requestContext"]
	authorizer, ok := reqContext.(map[string]interface{})["authorizer"]
	if !ok || authorizer == nil {
		return AuthProperties{}, errors.New("authorizer index missing in event")
	}

	claims, ok := authorizer.(map[string]interface{})["claims"]
	if !ok {
		return AuthProperties{}, errors.New("claims index missing in authorizer")
	}

	var authProps map[string]interface{}
	if value, ok := claims.(string); ok {
		err := json.Unmarshal([]byte(value), &authProps)
		if err != nil {
			logger.WithFields(logger.Fields{
				"error": err,
			}).Error("CognitoAuthProvider::ParseAuth() Could not parse claims as JSON")
			return AuthProperties{}, errors.New("could not parse claims as JSON")
		}
	} else {
		authProps = claims.(map[string]interface{})
	}
	return NewAuthProperties(authProps), nil
}
