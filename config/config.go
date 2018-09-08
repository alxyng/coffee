package config

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func newSession() (*session.Session, error) {
	if os.Getenv("AWS_SAM_LOCAL") == "true" {
		return session.NewSession(&aws.Config{
			Endpoint: aws.String("http://db:8000"),
			Region:   aws.String("eu-west-2"),
		})
	}

	return session.NewSession()
}

func CreateAWSSession() *session.Session {
	return session.Must(newSession())
}
