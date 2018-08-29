package s3

import (
	"errors"
	"github.com/matthisstenius/lambda-router/domain"
	"regexp"
	"strings"
)

const EventSource = "aws:s3"

// Route mapping for handler
type Route struct {
	Handler func(i *Input) domain.Response
}

// Routes mappings for S3 handlers
type Routes map[string]Route

// Router for S3 events
type Router struct {
	routes Routes
}

// NewRouter initializer
func NewRouter(routes Routes) *Router {
	return &Router{routes: routes}
}

// Route incoming event to corresponding handler
func (r *Router) Route(evt map[string]interface{}) (domain.Response, error) {
	record := evt["Records"].([]interface{})[0].(map[string]interface{})
	key := record["s3"].(map[string]interface{})["object"].(map[string]interface{})["key"].(string)

	re := regexp.MustCompile("[^/]+$")
	folder := re.ReplaceAllString(key, "")
	folder = "/" + strings.TrimSuffix(folder, "/")

	route, ok := r.routes[folder]
	if !ok {
		return nil, errors.New("handler func missing")
	}
	return route.Handler(NewInput(evt)), nil
}

// IsMatch for S3 event
func (r *Router) IsMatch(e map[string]interface{}) bool {
	if v, ok := e["Records"].([]interface{}); ok && len(v) > 0 {
		return v[0].(map[string]interface{})["eventSource"] == EventSource
	}
	return false
}
