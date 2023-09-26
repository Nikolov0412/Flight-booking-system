package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func initDynamoDB() (*dynamodb.DynamoDB, error) {
	awsConfig := &aws.Config{
		Region:   aws.String("us-east-1"),             // Change to your desired region.
		Endpoint: aws.String("http://localhost:4566"), // LocalStack endpoint
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	return dynamodb.New(sess), nil
}
