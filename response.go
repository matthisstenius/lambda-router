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
    encoded, _ := json.Marshal(body)

    log.Printf("Response body: %s", encoded)
    return &Response{
        StatusCode:      status,
        Body:            string(encoded),
        Headers:         map[string]string{},
        IsBase64Encoded: false,
    }
}

func NewErrorResponse(status int, error interface{}) *Response {
    encoded, _ := json.Marshal(map[string]interface{}{
        "error": error,
    })

    log.Printf("Error response: %s", encoded)
    return &Response{
        StatusCode: status,
        Body:       string(encoded),
        Headers: map[string]string{
            "Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key",
            "Access-Control-Allow-Methods": "*",
            "Access-Control-Allow-Origin":  "*",
        },
        IsBase64Encoded: false,
    }
}
