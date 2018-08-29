package mock

// Response mock
type Response struct {
	PayloadFn func() interface{}
}

// Payload mock implementation
func (r *Response) Payload() interface{} {
	return ""
}
