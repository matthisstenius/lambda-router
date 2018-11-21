package cognito

import (
	"github.com/matthisstenius/lambda-router/v2/http/cognito"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetParam(t *testing.T) {
	tests := []struct {
		Name   string
		Param  string
		Claims map[string]interface{}
		Out    interface{}
	}{
		{
			Name:   "it should succeed",
			Claims: map[string]interface{}{"custom:id": "1"},
			Param:  "custom:id",
			Out:    "1",
		},
		{
			Name:   "it should handle missing prop",
			Claims: map[string]interface{}{"custom:id": "1"},
			Param:  "missing",
			Out:    nil,
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// Given
			props := cognito.NewAuthProperties(td.Claims)

			// When
			prop := props.GetParam(td.Param)

			// Then
			assert.Equal(t, td.Out, prop)
		})
	}
}

func TestHasRole(t *testing.T) {
	tests := []struct {
		Name   string
		Role   string
		Claims map[string]interface{}
		Out    bool
	}{
		{
			Name:   "it should succeed on match",
			Claims: map[string]interface{}{"cognito:groups": "Admin"},
			Role:   "Admin",
			Out:    true,
		},
		{
			Name:   "it should handle mismatch",
			Claims: map[string]interface{}{"cognito:groups": "Admin"},
			Role:   "another",
			Out:    false,
		},
		{
			Name:   "it should handle missing cognito:groups in claims",
			Claims: map[string]interface{}{},
			Role:   "Admin",
			Out:    false,
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// Given
			props := cognito.NewAuthProperties(td.Claims)

			// When
			hasRole := props.HasRole(td.Role)

			// Then
			assert.Equal(t, td.Out, hasRole)
		})
	}
}
