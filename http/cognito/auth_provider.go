package cognito

import (
	"encoding/json"
	"errors"
	"github.com/matthisstenius/lambda-router/v2/domain"
	"github.com/matthisstenius/logger"
)

// AuthProvider Cognito implementation
type AuthProvider struct{}

// ParseAuth data from event
func (ap *AuthProvider) ParseAuth(evt map[string]interface{}) (domain.AuthProperties, error) {
	reqContext := evt["requestContext"]
	authorizer, ok := reqContext.(map[string]interface{})["authorizer"]
	if !ok || authorizer == nil {
		return nil, errors.New("authorizer index missing in event")
	}

	claims, ok := authorizer.(map[string]interface{})["claims"]
	if !ok {
		return nil, errors.New("claims index missing in authorizer")
	}

	var authProps domain.AuthProperties
	if value, ok := claims.(string); ok {
		err := json.Unmarshal([]byte(value), &authProps)
		if err != nil {
			logger.WithFields(logger.Fields{
				"error": err,
			}).Error("CognitoAuthProvider::ParseAuth() Could not parse claims as JSON")
			return nil, errors.New("could not parse claims as JSON")
		}
	} else {
		authProps = claims.(map[string]interface{})
	}
	return authProps, nil
}
