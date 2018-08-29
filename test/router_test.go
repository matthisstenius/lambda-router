package test

import (
	"testing"

	"errors"
	"github.com/matthisstenius/lambda-router"
	"github.com/matthisstenius/lambda-router/domain"
	"github.com/matthisstenius/lambda-router/mock"
	"github.com/stretchr/testify/assert"
)

var (
	httpRouterMock      *mock.Router
	dynamoRouterMock    *mock.Router
	s3RouterMock        *mock.Router
	scheduledRouterMock *mock.Router
	snsRouterMock       *mock.Router
)

func init() {
	httpRouterMock = new(mock.Router)
	dynamoRouterMock = new(mock.Router)
	s3RouterMock = new(mock.Router)
	scheduledRouterMock = new(mock.Router)
	snsRouterMock = new(mock.Router)
}

func TestStart(t *testing.T) {
	tests := []struct {
		Name        string
		Res         domain.Response
		IsHTTP      bool
		IsDynamo    bool
		IsS3        bool
		IsScheduled bool
		IsSNS       bool
		Error       error
	}{
		{
			Name:   "it should route http event",
			Res:    new(mock.Response),
			IsHTTP: true,
		},
		{
			Name:     "it should route dynamoDB stream event",
			Res:      new(mock.Response),
			IsDynamo: true,
		},
		{
			Name: "it should route S3 event",
			Res:  new(mock.Response),
			IsS3: true,
		},
		{
			Name:        "it should route scheduled event",
			Res:         new(mock.Response),
			IsScheduled: true,
		},
		{
			Name:  "it should route SNS event",
			Res:   new(mock.Response),
			IsSNS: true,
		},
		{
			Name:  "it should handle unknown event",
			Error: errors.New("unknown event"),
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// Given
			config := router.Config{
				HTTP:      httpRouterMock,
				DynamoDB:  dynamoRouterMock,
				S3:        s3RouterMock,
				Scheduled: scheduledRouterMock,
				SNS:       snsRouterMock,
			}

			httpRouterMock.IsMatchFn = func(evt map[string]interface{}) bool {
				return td.IsHTTP
			}

			httpRouterMock.DispatchFn = func(evt map[string]interface{}) (domain.Response, error) {
				return td.Res, td.Error
			}

			dynamoRouterMock.IsMatchFn = func(evt map[string]interface{}) bool {
				return td.IsDynamo
			}

			dynamoRouterMock.DispatchFn = func(evt map[string]interface{}) (domain.Response, error) {
				return td.Res, td.Error
			}

			s3RouterMock.IsMatchFn = func(evt map[string]interface{}) bool {
				return td.IsS3
			}

			s3RouterMock.DispatchFn = func(evt map[string]interface{}) (domain.Response, error) {
				return td.Res, td.Error
			}

			scheduledRouterMock.IsMatchFn = func(evt map[string]interface{}) bool {
				return td.IsScheduled
			}

			scheduledRouterMock.DispatchFn = func(evt map[string]interface{}) (domain.Response, error) {
				return td.Res, td.Error
			}

			snsRouterMock.IsMatchFn = func(evt map[string]interface{}) bool {
				return td.IsSNS
			}

			snsRouterMock.DispatchFn = func(evt map[string]interface{}) (domain.Response, error) {
				return td.Res, td.Error
			}

			// When
			event := router.NewEvent(&config)
			_, err := event.Handle(map[string]interface{}{})

			// Then
			assert.Equal(t, td.Error, err)
		})
	}
}
