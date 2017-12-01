package api

import (
	"encoding/json"
	"errors"
)

type Input struct {
	event map[string]interface{}
}

func (i Input) GetPathParam(param string) string {
	params, ok := i.event["pathParameters"]
	if !ok {
		return ""
	}
	return params.(map[string]string)[param]
}

func (i Input) GetQueryParam(param string) string {
	params, ok := i.event["queryStringParameters"]
	if !ok {
		return ""
	}
	return params.(map[string]string)[param]
}

func (i Input) PopulateBody(out interface{}) error {
	body, ok := i.event["body"]

	if !ok {
		return errors.New("missing request body")
	}

	err := json.Unmarshal([]byte(body.(string)), out)
	if err != nil {
		return errors.New("could not parse body as JSON")
	}

	return nil
}
