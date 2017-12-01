package api

import (
	"encoding/json"
	"errors"
	"github.com/gobuffalo/tags"
)

type Input struct {
	event map[string]interface{}
}

func (i Input) GetPathParam(param string) string {
	params, ok := i.event["pathParameters"]
	if !ok {
		return ""
	}
	value, ok := params.(map[string]interface{})[param]
	if !ok {
		return ""
	}
	return value.(string)
}

func (i Input) GetQueryParam(param string) string {
	params, ok := i.event["queryStringParameters"]
	if !ok {
		return ""
	}

	value, ok := params.(map[string]interface{})[param]
	if !ok {
		return ""
	}
	return value.(string)
}

func (i Input) PopulateBody(out interface{}) error {
	body, ok := i.event["body"]

	if !ok {
		return errors.New("missing request body")
	}

	err := json.Unmarshal([]byte(body.(string)), &out)
	if err != nil {
		return errors.New("could not parse body as JSON")
	}

	return nil
}

func (i Input) Body() string {
	return i.event["body"].(string)
}
