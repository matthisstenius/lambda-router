package api

import (
    "encoding/json"
    "log"
)

type Response struct {
    StatusCode      int               `json:"statusCode"`
    Body            interface{}       `json:"body"`
    Headers         map[string]string `json:"headers"`
    IsBase64Encoded bool              `json:"isBase64Encoded"`
}

func NewResponse(status int, body interface{}) *Response {
    encoded, _ := json.Marshal(map[string]interface{}{
        "code": status,
        "data": body,
    })

    log.Printf("Response: %s", encoded)
    return &Response{
        StatusCode:      status,
        Body:            string(encoded),
        Headers:         map[string]string{},
        IsBase64Encoded: false,
    }
}

func NewErrorResponse(status int, error interface{}) *Response {
    encoded, _ := json.Marshal(map[string]interface{}{
        "code":  status,
        "error": error,
    })

    log.Printf("Error response: %s", encoded)
    return &Response{
        StatusCode:      status,
        Body:            string(encoded),
        Headers:         map[string]string{},
        IsBase64Encoded: false,
    }
}
