package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func initDynamoDB() (*dynamodb.DynamoDB, error) {
	endpoint := os.Getenv("DYNAMODB_ENDPOINT") // Read the endpoint from an environment variable
	awsConfig := &aws.Config{
		Region:     aws.String("us-east-1"), // Change to your desired region.
		Endpoint:   aws.String(endpoint),    // LocalStack endpoint
		DisableSSL: aws.Bool(true),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	return dynamodb.New(sess), nil
}
