package schedule

import (
	"errors"
	"github.com/matthisstenius/lambda-router/domain"
	"github.com/matthisstenius/lambda-router/schedule"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoute(t *testing.T) {
	tests := []struct {
		Name     string
		Event    map[string]interface{}
		Schedule string
		Error    error
	}{
		{
			Name: "it should succeed",
			Event: map[string]interface{}{
				"resource": "test-schedule",
			},
			Schedule: "test-schedule",
		},
		{
			Name: "it should handle resource mismatch",
			Event: map[string]interface{}{
				"resource": "test-schedule",
			},
			Schedule: "other-schedule",
			Error:    errors.New("handler func missing"),
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// Given
			routes := schedule.Routes{
				td.Schedule: {
					Handler: func() domain.Response {
						return schedule.NewResponse("Success")
					},
				},
			}

			// When
			router := schedule.NewRouter(routes)
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
			Name:    "it should succeed",
			Event:   map[string]interface{}{"eventSource": schedule.EventSource},
			IsMatch: true,
		},
		{
			Name:    "it should handle none schedule event source",
			Event:   map[string]interface{}{"eventSource": "other:source"},
			IsMatch: false,
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// When
			router := schedule.NewRouter(schedule.Routes{})
			isMatch := router.IsMatch(td.Event)

			// Then
			assert.Equal(t, td.IsMatch, isMatch)
		})
	}
}
