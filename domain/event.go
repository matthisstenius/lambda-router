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

// Access DTO for roles and provider
type Access struct {
	Roles    []string
	Provider AccessProvider
}

// AccessProvider used for interpret incoming access parameters in event
type AccessProvider interface {
	Roles(evt map[string]interface{}) ([]string, error)
}
