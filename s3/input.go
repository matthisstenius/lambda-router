package s3

import "strings"

// Input for parsed S3 event
// TODO: Write tests
type Input struct {
	event map[string]interface{}
}

// NewInput initializer
func NewInput(e map[string]interface{}) *Input {
	return &Input{event: e}
}

// ObjectKeyPath extract full object key path
func (i *Input) ObjectKeyPath() string {
	record := i.event["Records"].([]interface{})[0].(map[string]interface{})
	return record["s3"].(map[string]interface{})["object"].(map[string]interface{})["key"].(string)
}

// ObjectKey extract object key from object key path
func (i *Input) ObjectKey() string {
	fragments := strings.Split(i.ObjectKeyPath(), "/")
	return fragments[len(fragments)-1]
}
