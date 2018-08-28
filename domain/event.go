package domain

type Response interface {
	Success(status int, body interface{}) Response
	Error(status int, error interface{}) Response
}

type Router interface {
	Route(evt map[string]interface{}) (Response, error)
	IsMatch(evt map[string]interface{}) bool
}
