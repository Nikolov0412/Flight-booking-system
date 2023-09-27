package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

type Airline struct {
	Name string `json:"Name"`
}

func ValidateAirlineName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) >= 6 {
		return errors.New("Airline name must have a length less than 6 characters")
	}
	return nil
}

func CreateAirline(name string, svc *dynamodb.DynamoDB) error {
	// Validate the airline name.
	if err := ValidateAirlineName(name); err != nil {
		return err
	}

	// Check if the airline name is already in use.
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("Airlines"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {
				S: aws.String(name),
			},
		},
		KeyConditionExpression: aws.String("Name = :name"),
	}
	result, err := svc.Query(queryInput)
	if err != nil {
		return err
	}

	if len(result.Items) > 0 {
		return errors.New("Airline name is not unique")
	}

	// Generate a unique ID for the airline.
	airlineID := uuid.New().String()

	// Create a new airline in DynamoDB with ID and Name.
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String("Airlines"),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(airlineID),
			},
			"Name": {
				S: aws.String(name),
			},
		},
	}
	_, err = svc.PutItem(putInput)
	if err != nil {
		return err
	}

	fmt.Printf("Created Airline: ID=%s, Name=%s\n", airlineID, name)
	return nil
}
