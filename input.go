package api

import (
    "encoding/json"
    "errors"
    "github.com/asaskevich/govalidator"
    log "github.com/sirupsen/logrus"
)

type Input struct {
    event map[string]interface{}
}

func (i *Input) GetPathParam(param string) string {
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

func (i *Input) GetQueryParam(param string) string {
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

func (i *Input) PopulateBody(out interface{}) error {
    body, ok := i.event["body"]
    if !ok || body == nil {
        log.WithField("body", body).Error("Input::PopulateBody() missing request body")
        return errors.New("missing request body")
    }

    if err := json.Unmarshal([]byte(body.(string)), &out); err != nil {
        log.WithField("error", err).Error("Input::PopulateBody() could not parse body as JSON")
        return errors.New("could not parse body as JSON")
    }

    if _, err := govalidator.ValidateStruct(out); err != nil {
        return err
    }
    return nil
}

// PopulateEventBody parse body in SNS event
func (i *Input) PopulateEventBody(out interface{}) error {
    record := i.event["Records"].([]map[string]interface{})[0]
    message := record["SNS"].(map[string]string)["Message"]

    if err := json.Unmarshal([]byte(message), &out); err != nil {
        return errors.New("could not parse SNS Message")
    }
    return nil
}

func (i *Input) Body() string {
    return i.event["body"].(string)
}

type CurrentUser struct {
    ID string `json:"id"`
}

func (i *Input) CurrentUser() *CurrentUser {
    reqContext := i.event["requestContext"]

    authorizer, ok := reqContext.(map[string]interface{})["authorizer"]
    if !ok || authorizer == nil {
        log.WithField("authorizer", authorizer).Fatal("Input::CurrentUser() authorizer index missing")
    }

    authData, ok := authorizer.(map[string]interface{})["authData"]
    if !ok || authorizer == nil {
        log.WithField("authData", authData).Fatal("Input::CurrentUser() authData index missing")
    }

    var parsed struct {
        Code int
        Data map[string]interface{}
    }

    err := json.Unmarshal([]byte(authData.(string)), &parsed)
    if err != nil {
        log.WithFields(log.Fields{
            "error": err,
            "data":  authData,
        }).Fatal("Input::CurrentUser() Could not parse authData")
    }
    return &CurrentUser{ID: parsed.Data["id"].(string)}
}
