package s3

import (
	"errors"
	"github.com/matthisstenius/lambda-router/v4/domain"
	"github.com/matthisstenius/lambda-router/v4/s3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoute(t *testing.T) {
	tests := []struct {
		Name   string
		Event  map[string]interface{}
		Folder string
		Error  error
	}{
		{
			Name: "it should succeed",
			Event: map[string]interface{}{
				"Records": []interface{}{
					map[string]interface{}{
						"s3": map[string]interface{}{
							"object": map[string]interface{}{
								"key": "folder/anotherFolder/object.test",
							},
						},
					},
				},
			},
			Folder: "/folder/anotherFolder",
		},
		{
			Name: "it should succeed with root folder",
			Event: map[string]interface{}{
				"Records": []interface{}{
					map[string]interface{}{
						"s3": map[string]interface{}{
							"object": map[string]interface{}{
								"key": "object.test",
							},
						},
					},
				},
			},
			Folder: "/",
		},
		{
			Name: "it should handle folder mismatch",
			Event: map[string]interface{}{
				"Records": []interface{}{
					map[string]interface{}{
						"s3": map[string]interface{}{
							"object": map[string]interface{}{
								"key": "folder/object.test",
							},
						},
					},
				},
			},
			Folder: "/otherFolder",
			Error:  errors.New("handler func missing"),
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// Given
			routes := s3.Routes{
				td.Folder: s3.Route{
					Handler: func(i *s3.Input) domain.Response {
						return s3.NewResponse("success")
					},
				},
			}

			// When
			router := s3.NewRouter(routes)
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
						"eventSource": s3.EventSource,
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
			router := s3.NewRouter(s3.Routes{})
			isMatch := router.IsMatch(td.Event)

			// Then
			assert.Equal(t, td.IsMatch, isMatch)
		})
	}
}
