package domain

// Response ...
type Response interface {
	Payload() interface{}
}

// Router ...
type Router interface {
	Route(evt map[string]interface{}) (Response, error)
	IsMatch(evt map[string]interface{}) bool
}

type Access struct {
	Roles    []string
	Provider AccessProvider
}

type AccessProvider interface {
	Roles(evt map[string]interface{}) ([]string, error)
}
