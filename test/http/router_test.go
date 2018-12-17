package http

import (
	"errors"
	"github.com/matthisstenius/lambda-router/v3/domain"
	"github.com/matthisstenius/lambda-router/v3/http"
	"github.com/matthisstenius/lambda-router/v3/mock"
	"github.com/stretchr/testify/assert"
	internalHTTP "net/http"
	"testing"
)

var (
	accessProviderMock *mock.AccessProvider
)

func init() {
	accessProviderMock = new(mock.AccessProvider)
}

func TestRoute(t *testing.T) {
	tests := []struct {
		Name                 string
		Event                map[string]interface{}
		Path                 string
		StatusCode           int
		MiddlewareStatusCode int
		HTTPMethod           string
		AccessProviderErr    error
		Roles                []string
		EventRoles           []string
		Middleware           []http.Middleware
		Error                error
	}{
		{
			Name: "it should succeed",
			Event: map[string]interface{}{
				"resource":   "/test/path",
				"httpMethod": internalHTTP.MethodGet,
			},
			Path:       "/test/path",
			HTTPMethod: internalHTTP.MethodGet,
			StatusCode: internalHTTP.StatusOK,
		},
		{
			Name: "it should succeed with access provider",
			Event: map[string]interface{}{
				"resource":   "/test/path",
				"httpMethod": internalHTTP.MethodGet,
			},
			Path:       "/test/path",
			Roles:      []string{"Admin"},
			EventRoles: []string{"Admin"},
			HTTPMethod: internalHTTP.MethodGet,
			StatusCode: internalHTTP.StatusOK,
		},
		{
			Name: "it should succeed with middleware that returns response",
			Event: map[string]interface{}{
				"resource":   "/test/path",
				"httpMethod": internalHTTP.MethodGet,
			},
			Middleware: []http.Middleware{
				func(i *http.Input) domain.Response {
					return http.NewErrorResponse(internalHTTP.StatusPaymentRequired, "Error")
				},
			},
			Path:                 "/test/path",
			HTTPMethod:           internalHTTP.MethodGet,
			StatusCode:           internalHTTP.StatusOK,
			MiddlewareStatusCode: internalHTTP.StatusPaymentRequired,
		},
		{
			Name: "it should succeed with middleware that returns nil",
			Event: map[string]interface{}{
				"resource":   "/test/path",
				"httpMethod": internalHTTP.MethodGet,
			},
			Middleware: []http.Middleware{
				func(i *http.Input) domain.Response {
					return nil
				},
			},
			Path:                 "/test/path",
			HTTPMethod:           internalHTTP.MethodGet,
			StatusCode:           internalHTTP.StatusOK,
			MiddlewareStatusCode: internalHTTP.StatusOK,
		},
		{
			Name: "it should handle access provider error",
			Event: map[string]interface{}{
				"resource":   "/protected/path",
				"httpMethod": internalHTTP.MethodGet,
			},
			Path:              "/protected/path",
			AccessProviderErr: errors.New("access provider error"),
			HTTPMethod:        internalHTTP.MethodGet,
			StatusCode:        internalHTTP.StatusForbidden,
		},
		{
			Name: "it should handle roles mismatch",
			Event: map[string]interface{}{
				"resource":   "/test/path",
				"httpMethod": internalHTTP.MethodGet,
			},
			Path:       "/test/path",
			Roles:      []string{"Admin"},
			EventRoles: []string{"Other"},
			HTTPMethod: internalHTTP.MethodGet,
			StatusCode: internalHTTP.StatusForbidden,
		},
		{
			Name: "it should succeed with path params",
			Event: map[string]interface{}{
				"resource":       "/test/1/other/2",
				"httpMethod":     internalHTTP.MethodGet,
				"pathParameters": map[string]interface{}{"id": "1", "otherID": "2"},
			},
			Path:       "/test/{id}/other/{otherID}",
			HTTPMethod: internalHTTP.MethodGet,
			StatusCode: internalHTTP.StatusOK,
		},
		{
			Name: "it should handle path mismatch",
			Event: map[string]interface{}{
				"resource":   "/test/path",
				"httpMethod": internalHTTP.MethodGet,
			},
			Path:       "/test/mismatch",
			HTTPMethod: internalHTTP.MethodGet,
			StatusCode: internalHTTP.StatusNotFound,
		},
		{
			Name: "it should handle method mismatch",
			Event: map[string]interface{}{
				"resource":   "/test/path",
				"httpMethod": internalHTTP.MethodPost,
			},
			Path:       "/test/path",
			HTTPMethod: internalHTTP.MethodGet,
			StatusCode: internalHTTP.StatusNotFound,
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// Given
			var routes http.Routes
			if len(td.Roles) > 0 {
				routes = http.Routes{
					td.Path: {
						td.HTTPMethod: http.Route{
							Handler: func(i *http.Input) domain.Response {
								return http.NewResponse(td.StatusCode, "")
							},
							Access: &domain.Access{
								Roles:    td.Roles,
								Provider: accessProviderMock,
							},
							Middleware: td.Middleware,
						},
					},
				}
			} else {
				routes = http.Routes{
					td.Path: {
						td.HTTPMethod: http.Route{
							Handler: func(i *http.Input) domain.Response {
								return http.NewResponse(td.StatusCode, "")
							},
							Middleware: td.Middleware,
						},
					},
				}
			}

			accessProviderMock.RolesFn = func(evt map[string]interface{}) ([]string, error) {
				return td.Roles, td.AccessProviderErr
			}

			// When
			router := http.NewRouter(routes)
			res, err := router.Route(td.Event)

			// Then
			assert.Equal(t, td.Error, err)

			if td.Error != nil {
				assert.Nil(t, res)
				return
			}
			if td.MiddlewareStatusCode != 0 {
				assert.Equal(t, td.MiddlewareStatusCode, res.Payload().(map[string]interface{})["statusCode"])
			} else {
				assert.Equal(t, td.StatusCode, res.Payload().(map[string]interface{})["statusCode"])
			}
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
				"httpMethod": internalHTTP.MethodPost,
			},
			IsMatch: true,
		},
		{
			Name:    "it should none match",
			Event:   map[string]interface{}{},
			IsMatch: false,
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// When
			router := http.NewRouter(http.Routes{})
			res := router.IsMatch(td.Event)

			// Then
			assert.Equal(t, td.IsMatch, res)
		})
	}
}
