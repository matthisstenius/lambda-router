package http

import (
	"errors"

	"github.com/matthisstenius/logger"
)

type AccessProvider interface {
	Roles(evt map[string]interface{}) ([]string, error)
}

type CognitoAccessProvider struct{}

func (ca *CognitoAccessProvider) Roles(evt map[string]interface{}) ([]string, error) {
	reqContext := evt["requestContext"]
	authorizer, ok := reqContext.(map[string]interface{})["authorizer"]
	if !ok || authorizer == nil {
		logger.WithFields(logger.Fields{}).Warning("CognitoAccessProvider::Roles() authorizer index missing in event")
		return nil, errors.New("authorizer index missing in event")
	}

	claims, ok := authorizer.(map[string]interface{})["claims"]
	if !ok {
		logger.WithFields(logger.Fields{}).Warning("CognitoAccessProvider::Roles() Claims missing in authorizer")
		return nil, errors.New("claims missing in authorizer")
	}
	roles, ok := claims.(map[string]interface{})["cognito:groups"]
	if !ok {
		logger.WithFields(logger.Fields{}).Info("CognitoAccessProvider::Roles() cognito:groups missing")
		return nil, errors.New("cognito:groups missing")
	}
	return []string{roles.(string)}, nil
}
