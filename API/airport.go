package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Airport struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
}

func ValidateAirportName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) != 3 || !isAlphabetic(name) {
		return errors.New("Airport name must be exactly 3 alphabetic characters")
	}

	return nil
}

// isAlphabetic checks if a string contains only alphabetic characters.
func isAlphabetic(s string) bool {
	return regexp.MustCompile("^[a-zA-Z]+$").MatchString(s)
}
func CreateAirport(name string, svc *dynamodb.DynamoDB) error {
	// Validate the airport name.
	if err := ValidateAirportName(name); err != nil {
		return err
	}

	// Check if the airport name is already in use.
	input := &dynamodb.QueryInput{
		TableName: aws.String("Airports"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {
				S: aws.String(name),
			},
		},
		KeyConditionExpression: aws.String("Code = :name"),
	}
	result, err := svc.Query(input)
	if err != nil {
		return err
	}

	if len(result.Items) > 0 {
		return errors.New("Airport name is not unique")
	}

	// Create a new airport in DynamoDB.
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String("Airports"),
		Item: map[string]*dynamodb.AttributeValue{
			"Code": {
				S: aws.String(name),
			},
		},
	}
	_, err = svc.PutItem(putInput)
	if err != nil {
		return err
	}

	fmt.Printf("Created Airport: %s\n", name)
	return nil
}
