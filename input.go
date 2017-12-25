package api

import (
	"encoding/json"
	"errors"
	"log"
)

type Input struct {
	event map[string]interface{}
}

func (i Input) GetPathParam(param string) string {
	params, ok := i.event["pathParameters"]
	if !ok || params == nil {
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
	if !ok || params == nil {
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

	if !ok || body == nil {
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

type CurrentUser struct {
	ID string `json:"id"`
}

func (i Input) CurrentUser() CurrentUser {
	reqContext := i.event["requestContext"]

	authorizer, ok := reqContext.(map[string]interface{})["authorizer"]

	if !ok || authorizer == nil {
		log.Fatal("Input::CurrentUser authorizer index missing")
	}

	authData, ok := authorizer.(map[string]interface{})["authData"]

	if !ok || authorizer == nil {
		log.Fatal("Input::CurrentUser authData index missing")
	}

	var parsed map[string]map[string]string
	err := json.Unmarshal([]byte(authData.(string)), &parsed)
	if err != nil {
		log.Fatal("Input::CurrentUser Could not parse authData")
	}

	return CurrentUser{ID: parsed["data"]["id"]}
}
