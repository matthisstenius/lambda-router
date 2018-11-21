package cognito

import (
	"errors"
	"github.com/matthisstenius/lambda-router/v2/http/cognito"
	"github.com/stretchr/testify/assert"
	"testing"
)

var accessProvider = new(cognito.AccessProvider)

func TestCognitoProvider(t *testing.T) {
	tests := []struct {
		Name  string
		Event map[string]interface{}
		Error error
		Roles []string
	}{
		{
			Name: "it should succeed by parsing roles",
			Event: map[string]interface{}{
				"requestContext": map[string]interface{}{
					"authorizer": map[string]interface{}{
						"claims": map[string]interface{}{
							"cognito:groups": "Admin",
						},
					},
				},
			},
			Roles: []string{"Admin"},
		},
		{
			Name: "it should handle missing claims",
			Event: map[string]interface{}{
				"requestContext": map[string]interface{}{
					"authorizer": map[string]interface{}{},
				},
			},
			Error: errors.New("claims index missing in authorizer"),
		},
		{
			Name: "it should handle missing cognito groups",
			Event: map[string]interface{}{
				"requestContext": map[string]interface{}{
					"authorizer": map[string]interface{}{
						"claims": map[string]interface{}{},
					},
				},
			},
			Error: errors.New("cognito:groups index missing in claims"),
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// When
			roles, err := accessProvider.ParseRoles(td.Event)

			// Then
			assert.Equal(t, td.Error, err)
			assert.Equal(t, td.Roles, roles)
		})
	}
}
