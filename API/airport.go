package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type Airport struct {
	Code string `json:"Code"`
}

func ValidateAirportCode(code string) error {
	code = strings.TrimSpace(code)
	if len(code) != 3 || !isAlphabetic(code) {
		return errors.New("Airport code must be exactly 3 alphabetic characters")
	}
	return nil
}

func CreateAirport(code string, svc *dynamodb.DynamoDB) error {
	// Validate the airport code.
	if err := ValidateAirportCode(code); err != nil {
		return err
	}

	// Check if the airport code is already in use.
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("Airports"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":code": {
				S: aws.String(code),
			},
		},
		KeyConditionExpression: aws.String("Code = :code"),
	}
	result, err := svc.Query(queryInput)
	if err != nil {
		return err
	}

	if len(result.Items) > 0 {
		return errors.New("Airport code is not unique")
	}

	// Generate a unique ID for the airport.
	airportID := uuid.New().String()

	// Create a new airport in DynamoDB with ID and Code.
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String("Airports"),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(airportID),
			},
			"Code": {
				S: aws.String(code),
			},
		},
	}
	_, err = svc.PutItem(putInput)
	if err != nil {
		return err
	}

	fmt.Printf("Created Airport: ID=%s, Code=%s\n", airportID, code)
	return nil
}

func isAlphabetic(s string) bool {
	return regexp.MustCompile("^[a-zA-Z]+$").MatchString(s)
}

func GetAirportByID(airportID string, svc *dynamodb.DynamoDB) (*Airport, error) {
	// Create a DynamoDB GetItem input.
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Airports"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(airportID),
			},
		},
	}

	// Get the item from DynamoDB.
	result, err := svc.GetItem(input)
	if err != nil {
		return nil, err
	}

	// Check if the item was found.
	if result.Item == nil {
		return nil, errors.New("Airport not found")
	}

	// Parse the retrieved data into an Airport struct.
	airport := &Airport{}
	err = dynamodbattribute.UnmarshalMap(result.Item, airport)
	if err != nil {
		return nil, err
	}

	return airport, nil
}
