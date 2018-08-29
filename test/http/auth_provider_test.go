package http

import (
	"errors"
	"github.com/matthisstenius/lambda-router/domain"
	"github.com/matthisstenius/lambda-router/http/cognito"
	"github.com/stretchr/testify/assert"
	"testing"
)

var authProvider = new(cognito.AuthProvider)

func TestCognitoAuthProvider(t *testing.T) {
	tests := []struct {
		Name      string
		Event     map[string]interface{}
		AuthProps domain.AuthProperties
		Error     error
	}{
		{
			Name: "it should succeed",
			Event: map[string]interface{}{
				"requestContext": map[string]interface{}{
					"authorizer": map[string]interface{}{
						"claims": map[string]interface{}{
							"custom:id": "12345",
						},
					},
				},
			},
			AuthProps: domain.AuthProperties{"custom:id": "12345"},
		},
		{
			Name: "it should succeed with claims as JSON",
			Event: map[string]interface{}{
				"requestContext": map[string]interface{}{
					"authorizer": map[string]interface{}{
						"claims": `{"custom:id": "12345"}`,
					},
				},
			},
			AuthProps: domain.AuthProperties{"custom:id": "12345"},
		},
		{
			Name: "it should handle missing claims index",
			Event: map[string]interface{}{
				"requestContext": map[string]interface{}{
					"authorizer": map[string]interface{}{},
				},
			},
			Error: errors.New("claims index missing in authorizer"),
		},
		{
			Name: "it should handle invalid claims JSON",
			Event: map[string]interface{}{
				"requestContext": map[string]interface{}{
					"authorizer": map[string]interface{}{
						"claims": `{"custom:id": invalid"}`,
					},
				},
			},
			Error: errors.New("could not parse claims as JSON"),
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// When
			props, err := authProvider.ParseAuth(td.Event)

			// Then
			assert.Equal(t, td.Error, err)
			assert.Equal(t, td.AuthProps, props)
		})
	}
}
