package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/matthisstenius/lambda-router/v3/domain"
)

type Middleware func(i *Input) domain.Response

// Route mapping for handler and optional access
type Route struct {
	Handler    func(i *Input) domain.Response
	Access     *domain.Access
	Middleware []Middleware
}

// Routes mappings for HTTP handlers
type Routes map[string]map[string]Route

// Router for HTTP events
type Router struct {
	routes     Routes
	middleware []Middleware
}

// NewRouter initializer
func NewRouter(routes Routes, middleware []Middleware) *Router {
	return &Router{routes: routes, middleware: middleware}
}

// Dispatch incoming event to corresponding handler
func (r *Router) Route(evt map[string]interface{}) (domain.Response, error) {
	pathParams, ok := evt["pathParameters"]
	resource := evt["resource"].(string)
	method := evt["httpMethod"].(string)

	if ok && pathParams != nil {
		for k, v := range pathParams.(map[string]interface{}) {
			resource = strings.Replace(resource, v.(string), fmt.Sprintf("{%s}", k), 1)
		}
	}

	route, ok := r.routes[resource][method]
	if !ok {
		return NewErrorResponse(http.StatusNotFound, "No matching handler found"), nil
	}

	i := NewInput(evt)
	if !r.hasAccess(route.Access, i) {
		return NewErrorResponse(http.StatusForbidden, "Access denied"), nil
	}

	for _, m := range route.Middleware {
		if res := m(i); res != nil {
			return res, nil
		}
	}
	for _, m := range r.middleware {
		if res := m(i); res != nil {
			return res, nil
		}
	}
	return route.Handler(i), nil
}

// IsMatch for HTTP event
func (r *Router) IsMatch(e map[string]interface{}) bool {
	if _, ok := e["httpMethod"]; ok {
		return true
	}
	return false
}

func (r *Router) hasAccess(access *domain.Access, i *Input) bool {
	if access == nil {
		return true
	}

	claims, err := i.Auth()
	if err != nil {
		return false
	}

	reqRole := claims.Get(access.Key)
	if reqRole == nil {
		return false
	}

	match := false
	for _, r := range access.Roles {
		if reqRole == r {
			match = true
		}
	}
	return match
}
