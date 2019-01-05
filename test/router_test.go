package test

import (
	"testing"

	"errors"
	"github.com/matthisstenius/lambda-router/v4"
	"github.com/matthisstenius/lambda-router/v4/domain"
	"github.com/matthisstenius/lambda-router/v4/mock"
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
		Name            string
		Res             domain.Response
		HTTPRouter      *mock.Router
		DynamoRouter    *mock.Router
		S3Router        *mock.Router
		SNSRouter       *mock.Router
		ScheduledRouter *mock.Router
		IsHTTP          bool
		IsDynamo        bool
		IsS3            bool
		IsScheduled     bool
		IsSNS           bool
		Error           error
	}{
		{
			Name:       "it should route http event",
			Res:        new(mock.Response),
			HTTPRouter: httpRouterMock,
			IsHTTP:     true,
		},
		{
			Name:         "it should route dynamoDB stream event",
			Res:          new(mock.Response),
			DynamoRouter: dynamoRouterMock,
			IsDynamo:     true,
		},
		{
			Name:     "it should route S3 event",
			Res:      new(mock.Response),
			S3Router: s3RouterMock,
			IsS3:     true,
		},
		{
			Name:            "it should route scheduled event",
			Res:             new(mock.Response),
			ScheduledRouter: scheduledRouterMock,
			IsScheduled:     true,
		},
		{
			Name:      "it should route SNS event",
			Res:       new(mock.Response),
			SNSRouter: snsRouterMock,
			IsSNS:     true,
		},
		{
			Name:  "it should handle unknown event",
			Error: errors.New("unknown event"),
		},
		{
			Name:  "it should handle missing router in config",
			Error: errors.New("unknown event"),
		},
	}

	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			// Given
			config := router.Config{}

			if td.HTTPRouter != nil {
				td.HTTPRouter.IsMatchFn = func(evt map[string]interface{}) bool {
					return td.IsHTTP
				}
				td.HTTPRouter.DispatchFn = func(evt map[string]interface{}) (domain.Response, error) {
					return td.Res, td.Error
				}
				config.HTTP = td.HTTPRouter
			}
			if td.DynamoRouter != nil {
				td.DynamoRouter.IsMatchFn = func(evt map[string]interface{}) bool {
					return td.IsDynamo
				}

				td.DynamoRouter.DispatchFn = func(evt map[string]interface{}) (domain.Response, error) {
					return td.Res, td.Error
				}
				config.DynamoDB = td.DynamoRouter
			}
			if td.S3Router != nil {
				td.S3Router.IsMatchFn = func(evt map[string]interface{}) bool {
					return td.IsS3
				}

				td.S3Router.DispatchFn = func(evt map[string]interface{}) (domain.Response, error) {
					return td.Res, td.Error
				}
				config.S3 = td.S3Router
			}
			if td.ScheduledRouter != nil {
				td.ScheduledRouter.IsMatchFn = func(evt map[string]interface{}) bool {
					return td.IsScheduled
				}

				td.ScheduledRouter.DispatchFn = func(evt map[string]interface{}) (domain.Response, error) {
					return td.Res, td.Error
				}
				config.Scheduled = td.ScheduledRouter
			}
			if td.SNSRouter != nil {
				td.SNSRouter.IsMatchFn = func(evt map[string]interface{}) bool {
					return td.IsSNS
				}

				td.SNSRouter.DispatchFn = func(evt map[string]interface{}) (domain.Response, error) {
					return td.Res, td.Error
				}
				config.SNS = td.SNSRouter
			}

			// When
			event := router.NewEvent(&config)
			_, err := event.Handle(map[string]interface{}{})

			// Then
			assert.Equal(t, td.Error, err)
		})
	}
}
