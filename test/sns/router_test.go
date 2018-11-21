package sns

import (
	"errors"
	"github.com/matthisstenius/lambda-router/v3/domain"
	"github.com/matthisstenius/lambda-router/v3/sns"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoute(t *testing.T) {
	tests := []struct {
		Name     string
		Event    map[string]interface{}
		TopicARN string
		Error    error
	}{
		{
			Name: "it should succeed",
			Event: map[string]interface{}{
				"Records": []interface{}{
					map[string]interface{}{
						"Sns": map[string]interface{}{
							"TopicArn": "test:topic:arn",
						},
					},
				},
			},
			TopicARN: "test:topic:arn",
		},
		{
			Name: "it should handle topic mismatch",
			Event: map[string]interface{}{
				"Records": []interface{}{
					map[string]interface{}{
						"Sns": map[string]interface{}{
							"TopicArn": "test:topic:arn",
						},
					},
				},
			},
			TopicARN: "test:other:arn",
			Error:    errors.New("handler func missing"),
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// Given
			routes := sns.Routes{
				td.TopicARN: {
					Handler: func(i *sns.Input) domain.Response {
						return sns.NewResponse("Success")
					},
				},
			}

			// When
			router := sns.NewRouter(routes)
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
						"EventSource": sns.EventSource,
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
			router := sns.NewRouter(sns.Routes{})
			isMatch := router.IsMatch(td.Event)

			// Then
			assert.Equal(t, td.IsMatch, isMatch)
		})
	}
}
