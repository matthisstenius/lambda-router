package http

import (
	"errors"
	"github.com/matthisstenius/lambda-router/v4/domain"
	"github.com/matthisstenius/lambda-router/v4/http"
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
		{
			Name: "it should handle none string value by encode to JSON",
			Event: map[string]interface{}{
				"pathParameters": map[string]interface{}{
					"testParam": []interface{}{"Test value"},
				},
			},
			Param: "testParam",
			Value: `["Test value"]`,
		},
		{
			Name:  "it should handle missing path parameters",
			Event: map[string]interface{}{},
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

func TestHasPathParam(t *testing.T) {
	tests := []struct {
		Name  string
		Event map[string]interface{}
		Param string
		Out   bool
	}{
		{
			Name: "it should succeed",
			Event: map[string]interface{}{
				"pathParameters": map[string]interface{}{
					"testParam": "test value",
				},
			},
			Param: "testParam",
			Out:   true,
		},
		{
			Name: "it should handle missing param",
			Event: map[string]interface{}{
				"pathParameters": map[string]interface{}{
					"testParam": "test value",
				},
			},
			Param: "otherParam",
			Out:   false,
		},
		{
			Name:  "it should handle missing path parameters",
			Event: map[string]interface{}{},
			Param: "otherParam",
			Out:   false,
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// Given, When
			input := http.NewInput(td.Event)
			match := input.HasPathParam(td.Param)

			// Then
			assert.Equal(t, td.Out, match)
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
		{
			Name: "it should handle none string value by encode to JSON",
			Event: map[string]interface{}{
				"queryStringParameters": map[string]interface{}{
					"testParam": []interface{}{"Test value"},
				},
			},
			Param: "testParam",
			Value: `["Test value"]`,
		},
		{
			Name:  "it should handle missing query string parameters",
			Event: map[string]interface{}{},
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

func TestGetHeader(t *testing.T) {
	tests := []struct {
		Name   string
		Event  map[string]interface{}
		Header string
		Out    interface{}
	}{
		{
			Name: "it should succeed",
			Event: map[string]interface{}{
				"headers": map[string]interface{}{
					"Test-Header": "test value",
				},
			},
			Header: "Test-Header",
			Out:    "test value",
		},
		{
			Name: "it should handle missing header",
			Event: map[string]interface{}{
				"headers": map[string]interface{}{
					"Test-Header": "test value",
				},
			},
			Header: "Other-Header",
			Out:    "",
		},
		{
			Name:   "it should handle missing header",
			Event:  map[string]interface{}{},
			Header: "Other-Header",
			Out:    "",
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// Given
			input := http.NewInput(td.Event)

			// When
			out := input.GetHeader(td.Header)

			// Then
			assert.Equal(t, td.Out, out)
		})
	}
}

func TestHasQueryParam(t *testing.T) {
	tests := []struct {
		Name  string
		Event map[string]interface{}
		Param string
		Out   bool
	}{
		{
			Name: "it should succeed",
			Event: map[string]interface{}{
				"queryStringParameters": map[string]interface{}{
					"testParam": "test value",
				},
			},
			Param: "testParam",
			Out:   true,
		},
		{
			Name:  "it should handle missing query string parameters",
			Event: map[string]interface{}{},
			Param: "otherParam",
			Out:   false,
		},
		{
			Name: "it should handle missing param",
			Event: map[string]interface{}{
				"queryStringParameters": map[string]interface{}{
					"testParam": "test value",
				},
			},
			Param: "otherParam",
			Out:   false,
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// Given, When
			input := http.NewInput(td.Event)
			match := input.HasQueryParam(td.Param)

			// Then
			assert.Equal(t, td.Out, match)
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
			Name: "it should succeed for parsed JSON",
			Event: map[string]interface{}{
				"queryStringParameters": map[string]interface{}{
					"testParam": []interface{}{"Test value"},
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

func TestRawBody(t *testing.T) {
	tests := []struct {
		Name  string
		Event map[string]interface{}
		Out   []byte
	}{
		{
			Name: "it should succeed",
			Event: map[string]interface{}{
				"body": `{"message": "hello, world"}`,
			},
			Out: []byte(`{"message": "hello, world"}`),
		},
		{
			Name:  "it should handle missing body",
			Event: map[string]interface{}{},
			Out:   []byte(""),
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// Given
			input := http.NewInput(td.Event)

			// When
			out := input.RawBody()

			// Then
			assert.Equal(t, td.Out, out)
		})
	}
}

func TestAuth(t *testing.T) {
	tests := []struct {
		Name  string
		Event map[string]interface{}
		Out   *domain.AuthClaims
		Error error
	}{
		{
			Name: "it should succeed",
			Event: map[string]interface{}{
				"requestContext": map[string]interface{}{
					"authorizer": map[string]interface{}{
						"claims": map[string]interface{}{
							"id": "12345",
						},
					},
				},
			},
			Out: domain.NewAuthClaims(map[string]interface{}{"id": "12345"}),
		},
		{
			Name: "it should succeed with claims as JSON",
			Event: map[string]interface{}{
				"requestContext": map[string]interface{}{
					"authorizer": map[string]interface{}{
						"claims": `{"id": "12345"}`,
					},
				},
			},
			Out: domain.NewAuthClaims(map[string]interface{}{"id": "12345"}),
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
			// Given
			input := http.NewInput(td.Event)

			// When
			claims, err := input.Auth()

			// Then
			assert.Equal(t, td.Error, err)
			assert.Equal(t, td.Out, claims)
		})
	}
}
