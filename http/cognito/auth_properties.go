package cognito

// AuthProperties cognito implementation
type AuthProperties struct {
	props map[string]interface{}
}

// NewAuthProperties initializer
func NewAuthProperties(props map[string]interface{}) *AuthProperties {
	return &AuthProperties{props: props}
}

// GetParam from claims
func (ap *AuthProperties) GetParam(key string) interface{} {
	if _, ok := ap.props[key]; !ok {
		return nil
	}
	return ap.props[key]
}

// HasRole check current role in claims
func (ap *AuthProperties) HasRole(role string) bool {
	if v, ok := ap.props["cognito:groups"]; ok {
		return v == role
	}
	return false
}
