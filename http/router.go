package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/matthisstenius/lambda-router/v2/domain"
)

// Route mapping for handler and optional access
type Route struct {
	Handler func(i *Input) domain.Response
	Access  *domain.Access
}

// Routes mappings for HTTP handlers
type Routes map[string]map[string]Route

// Router for HTTP events
type Router struct {
	routes Routes
}

// NewRouter initializer
func NewRouter(routes Routes) *Router {
	return &Router{routes: routes}
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

	if !r.hasAccess(route.Access, evt) {
		return NewErrorResponse(http.StatusForbidden, "Access denied"), nil
	}

	return route.Handler(NewInput(evt)), nil
}

// IsMatch for HTTP event
func (r *Router) IsMatch(e map[string]interface{}) bool {
	if _, ok := e["httpMethod"]; ok {
		return true
	}
	return false
}

func (r *Router) hasAccess(access *domain.Access, evt map[string]interface{}) bool {
	if access == nil {
		return true
	}

	reqRoles, err := access.Provider.ParseRoles(evt)
	if err != nil {
		return false
	}

	var roleMatch bool
	for _, re := range reqRoles {
		for _, r := range access.Roles {
			if re == r {
				roleMatch = true
			}
		}
	}

	if !roleMatch {
		return false
	}
	return true
}
