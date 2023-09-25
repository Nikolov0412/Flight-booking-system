package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func initDynamoDB() (*dynamodb.DynamoDB, error) {
	awsConfig := &aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("00000000", "00000000", ""),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	return dynamodb.New(sess), nil
}
