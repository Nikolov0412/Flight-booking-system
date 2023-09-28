package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

func GetAirlineByID(airlineID string, svc *dynamodb.DynamoDB) (*Airline, error) {
	// Create a DynamoDB GetItem input.
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Airlines"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(airlineID),
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
		return nil, errors.New("Airline not found")
	}

	// Parse the retrieved data into an Airline struct.
	airline := &Airline{}
	err = dynamodbattribute.UnmarshalMap(result.Item, airline)
	if err != nil {
		return nil, err
	}

	return airline, nil
}

func GetAllAirlines(svc *dynamodb.DynamoDB) ([]*Airline, error) {
	// Create a DynamoDB Scan input to retrieve all items from the Airlines table.
	input := &dynamodb.ScanInput{
		TableName: aws.String("Airlines"),
	}

	// Perform the scan operation.
	result, err := svc.Scan(input)
	if err != nil {
		return nil, err
	}

	// Initialize a slice to hold the retrieved airlines.
	airlines := []*Airline{}

	// Iterate through the scan results and parse each item into an Airline struct.
	for _, item := range result.Items {
		airline := &Airline{}
		if err := dynamodbattribute.UnmarshalMap(item, airline); err != nil {
			return nil, err
		}
		airlines = append(airlines, airline)
	}

	return airlines, nil
}
