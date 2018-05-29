package router

import (
	"errors"
	"lambda-router/dynamodb"
	"runtime/debug"

	"lambda-router/http"

	"lambda-router/schedule"

	"lambda-router/s3"

	"lambda-router/sns"

	"github.com/matthisstenius/logger"
)

// Router that can be used to route different AWS Lambda events
// to corresponding handlers
type Router struct {
	config *Config
}

// Config for routing event handlers
type Config struct {
	HTTP      http.Routes
	Scheduled schedule.Routes
	DynamoDB  dynamodb.Routes
	S3        s3.Routes
	SNS       sns.Routes
}

// NewHandler initialization for Router
func NewRouter(config *Config) *Router {
	return &Router{config: config}
}

// Start correct handler based in mapped handler routes and incoming event
func (r *Router) Start(event interface{}) (interface{}, error) {
	e := event.(map[string]interface{})
	logger.WithFields(logger.Fields{
		"event": event,
	}).Info("Incoming event")
	defer r.logPanic()

	var response interface{}
	var err error
	switch true {
	case http.IsMatch(e):
		router := http.NewRouter(e, r.config.HTTP)
		response, err = router.Dispatch()
		break
	case schedule.IsMatch(e):
		router := schedule.NewRouter(e, r.config.Scheduled)
		response, err = router.Dispatch()
		break
	case dynamodb.IsMatch(e):
		router := dynamodb.NewRouter(e, r.config.DynamoDB)
		response, err = router.Dispatch()
		break
	case s3.IsMatch(e):
		router := s3.NewRouter(e, r.config.S3)
		response, err = router.Dispatch()
		break
	case sns.IsMatch(e):
		router := sns.NewRouter(e, r.config.SNS)
		response, err = router.Dispatch()
		break
	default:
		response, err = nil, errors.New("unknown event")
	}
	return response, err
}

func (r *Router) logPanic() {
	if r := recover(); r != nil {
		logger.WithFields(logger.Fields{
			"error": r,
			"stack": string(debug.Stack()),
		}).Error("Unexpected panic")
	}
}
