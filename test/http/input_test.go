package http

import (
	"errors"
	"github.com/matthisstenius/lambda-router/v2/domain"
	"github.com/matthisstenius/lambda-router/v2/http"
	"github.com/matthisstenius/lambda-router/v2/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPathParam(t *testing.T) {
	tests := []struct {
		Name  string
		Event map[string]interface{}
		Param string
		Value string
	}{
		{
			Name: "it should succeed",
			Event: map[string]interface{}{
				"pathParameters": map[string]interface{}{
					"testParam": "test value",
				},
			},
			Param: "testParam",
			Value: "test value",
		},
		{
			Name: "it should handle missing param",
			Event: map[string]interface{}{
				"pathParameters": map[string]interface{}{
					"testParam": "test value",
				},
			},
			Param: "otherParam",
			Value: "",
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// Given

			// When
			input := http.NewInput(td.Event)
			value := input.GetPathParam(td.Param)

			// Then
			assert.Equal(t, td.Value, value)
		})
	}
}

func TestGetQueryParam(t *testing.T) {
	tests := []struct {
		Name  string
		Event map[string]interface{}
		Param string
		Value interface{}
	}{
		{
			Name: "it should succeed",
			Event: map[string]interface{}{
				"queryStringParameters": map[string]interface{}{
					"testParam": "test value",
				},
			},
			Param: "testParam",
			Value: "test value",
		},
		{
			Name: "it should handle missing param",
			Event: map[string]interface{}{
				"queryStringParameters": map[string]interface{}{
					"testParam": "test value",
				},
			},
			Param: "otherParam",
			Value: "",
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// Given

			// When
			input := http.NewInput(td.Event)
			value := input.GetQueryParam(td.Param)

			// Then
			assert.Equal(t, td.Value, value)
		})
	}
}

func TestParseQueryParam(t *testing.T) {
	tests := []struct {
		Name          string
		Event         map[string]interface{}
		Param         string
		Out           []string
		ExpectedValue []string
		Error         error
	}{
		{
			Name: "it should succeed",
			Event: map[string]interface{}{
				"queryStringParameters": map[string]interface{}{
					"testParam": `["Test value"]`,
				},
			},
			Param:         "testParam",
			ExpectedValue: []string{"Test value"},
		},
		{
			Name: "it should handle invalid JSON",
			Event: map[string]interface{}{
				"queryStringParameters": map[string]interface{}{
					"testParam": `Test value`,
				},
			},
			Param: "testParam",
			Error: errors.New("could not parse param as JSON"),
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			input := http.NewInput(td.Event)
			err := input.ParseQueryParam(td.Param, &td.Out)

			// Then
			assert.Equal(t, td.Error, err)
			assert.Equal(t, td.ExpectedValue, td.Out)
		})
	}
}

func TestParseBody(t *testing.T) {
	tests := []struct {
		Name  string
		Event map[string]interface{}
		Body  map[string]string
		Error error
	}{
		{
			Name: "it should succeed",
			Event: map[string]interface{}{
				"body": `{"message": "hello, world"}`,
			},
			Body: map[string]string{"message": "hello, world"},
		},
		{
			Name:  "it should handle missing body",
			Event: map[string]interface{}{},
			Error: errors.New("missing request body"),
		},
		{
			Name: "it should handle invalid json",
			Event: map[string]interface{}{
				"body": `{"message: "invalid"}`,
			},
			Error: errors.New("could not parse body as JSON"),
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// Given

			// When
			input := http.NewInput(td.Event)
			var body map[string]string
			err := input.ParseBody(&body)

			// Then
			assert.Equal(t, td.Error, err)
			assert.Equal(t, td.Body, body)
		})
	}
}

func TestAuth(t *testing.T) {
	tests := []struct {
		Name         string
		AuthProvider domain.AuthProvider
		Error        error
	}{
		{
			Name:         "it should succeed",
			AuthProvider: new(mock.AuthProvider),
		},
		{
			Name:  "it should handle missing auth provider",
			Error: errors.New("given auth provider is nil"),
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// When
			input := http.NewInput(map[string]interface{}{})
			_, err := input.Auth(td.AuthProvider)

			// Then
			assert.Equal(t, td.Error, err)
		})
	}
}
