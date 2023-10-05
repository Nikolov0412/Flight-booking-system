package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func doesTableExist(tableName string, svc *dynamodb.DynamoDB) bool {
	// Describe the table to check if it exists.
	_, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	return err == nil
}
