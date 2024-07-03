package utils

import (
	"fmt"
	"order-service/internal/constants"
	"os"
	"slices"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func ReadConfigFile(path string) string {
	// Valid envs: local, dev, uat, prod
	envs := []string{"local", "dev", "uat", "prod"}

	env := os.Getenv("ENV")
	if env == constants.EmptyString || !slices.Contains(envs, env) {
		fmt.Println("No valid env provided, defaulting to local")
		env = "local"
	}

	filePath := fmt.Sprintf("%s/config_%s.yaml", path, env)
	return filePath
}

func CreateAWSSession(region string) *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewEnvCredentials(),
		},
	}))
}
