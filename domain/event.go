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
	Roles []string
	Key   string
}

// AuthClaims for current authenticated user
type AuthClaims struct {
	claims map[string]interface{}
}

// NewAuthClaims initializer
func NewAuthClaims(claims map[string]interface{}) *AuthClaims {
	return &AuthClaims{claims: claims}
}

// Get claim by key
func (ac *AuthClaims) Get(key string) interface{} {
	if _, ok := ac.claims[key]; !ok {
		return nil
	}
	return ac.claims[key]
}
