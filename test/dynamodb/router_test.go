package dynamodb

import (
	"errors"
	"github.com/matthisstenius/lambda-router/domain"
	"github.com/matthisstenius/lambda-router/dynamodb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoute(t *testing.T) {
	tests := []struct {
		Name   string
		Event  map[string]interface{}
		Stream string
		Error  error
	}{
		{
			Name: "it should succeed",
			Event: map[string]interface{}{
				"Records": []interface{}{
					map[string]interface{}{
						"eventSourceARN": "test-stream",
					},
				},
			},
			Stream: "test-stream",
		},
		{
			Name: "it should handle stream mismatch",
			Event: map[string]interface{}{
				"Records": []interface{}{
					map[string]interface{}{
						"eventSourceARN": "test-stream",
					},
				},
			},
			Stream: "other-stream",
			Error:  errors.New("handler func missing"),
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// Given
			routes := dynamodb.Routes{
				td.Stream: dynamodb.Route{
					Handler: func(i *dynamodb.Input) domain.Response {
						return dynamodb.NewResponse("success")
					},
				},
			}

			// When
			router := dynamodb.NewRouter(routes)
			_, err := router.Route(td.Event)

			// Then
			assert.Equal(t, td.Error, err)
		})
	}
}

func TestIsMatch(t *testing.T) {
	tests := []struct {
		Name    string
		Event   map[string]interface{}
		IsMatch bool
	}{
		{
			Name: "it should succeed",
			Event: map[string]interface{}{
				"Records": []interface{}{
					map[string]interface{}{
						"eventSource": dynamodb.EventSource,
					},
				},
			},
			IsMatch: true,
		},
		{
			Name:    "it should handle missing records",
			Event:   map[string]interface{}{},
			IsMatch: false,
		},
		{
			Name: "it should handle empty records",
			Event: map[string]interface{}{
				"Records": []interface{}{},
			},
			IsMatch: false,
		},
		{
			Name: "it should handle none dynamo event source",
			Event: map[string]interface{}{
				"Records": []interface{}{
					map[string]interface{}{
						"eventSource": "other:source",
					},
				},
			},
			IsMatch: false,
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// When
			router := dynamodb.NewRouter(dynamodb.Routes{})
			isMatch := router.IsMatch(td.Event)

			// Then
			assert.Equal(t, td.IsMatch, isMatch)
		})
	}
}
