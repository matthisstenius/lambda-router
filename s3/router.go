package s3

import (
	"errors"
	"regexp"
	"strings"
)

const eventSourceS3 = "aws:s3"

// Routes mappings for S3 handlers
type Routes map[string]func(i *Input) *Response

// Router for S3 events
type Router struct {
	event  map[string]interface{}
	routes Routes
}

// NewRouter initializer
func NewRouter(e map[string]interface{}, routes Routes) *Router {
	return &Router{event: e, routes: routes}
}

// Dispatch incoming event to corresponding handler
func (r *Router) Dispatch() (*Response, error) {
	record := r.event["Records"].([]interface{})[0].(map[string]interface{})
	key := record["s3"].(map[string]interface{})["object"].(map[string]interface{})["key"].(string)

	re := regexp.MustCompile("[^/]+$")
	folder := re.ReplaceAllString(key, "")
	if folder == "" {
		// If object is in root we want to look to /
		folder = "/"
	} else {
		folder = strings.TrimSuffix(folder, "/")
	}

	handler, ok := r.routes[folder]
	if !ok {
		return nil, errors.New("handler func missing")
	}
	return &*handler(&Input{event: r.event}), nil
}

// IsMatch for S3 event
func IsMatch(e map[string]interface{}) bool {
	if v, ok := e["Records"].([]interface{}); ok && len(v) > 0 {
		return v[0].(map[string]interface{})["eventSource"] == eventSourceS3
	}
	return false
}
