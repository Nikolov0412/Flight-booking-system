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
	ID   string `json:"id"`
	Code string `json:"code"`
}

func ValidateAirlineCode(code string) error {
	code = strings.TrimSpace(code)
	if len(code) >= 6 {
		return errors.New("Airline code must have a length less than 6 characters")
	}
	return nil
}

func CreateAirline(airline Airline, svc *dynamodb.DynamoDB) error {
	// Validate the airline code.
	if err := ValidateAirlineCode(airline.Code); err != nil {
		return err
	}
	if !doesTableExist("Airlines", svc) {
		if err := createAirlinesTable(svc); err != nil {
			fmt.Printf("Error creating Airlines table: %v\n", err)
		}
	}
	// Check if the airline code is already in use.
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("Airlines"),
		IndexName: aws.String("CodeIndex"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":code": {
				S: aws.String(airline.Code),
			},
		},
		KeyConditionExpression: aws.String("Code = :code"),
	}
	result, err := svc.Query(queryInput)
	if err != nil {
		return err
	}

	if len(result.Items) > 0 {
		return errors.New("Airline code is not unique")
	}

	airlineID := uuid.New().String()
	// Create a new airline in DynamoDB with ID and Code.
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String("Airlines"),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(airlineID),
			},
			"Code": {
				S: aws.String(airline.Code),
			},
		},
	}
	_, err = svc.PutItem(putInput)
	if err != nil {
		return err
	}

	fmt.Printf("Created Airline: ID=%s, Code=%s\n", airlineID, airline.Code)
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

func createAirlinesTable(svc *dynamodb.DynamoDB) error {
	// Define the parameters for creating the "Airlines" table.
	params := &dynamodb.CreateTableInput{
		TableName: aws.String("Airlines"),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Code"),
				KeyType:       aws.String("RANGE"),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Code"),
				AttributeType: aws.String("S"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("CodeIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("Code"),
						KeyType:       aws.String("HASH"), // Secondary index key
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			},
		},
	}

	// Create the "Airlines" table.
	_, err := svc.CreateTable(params)
	if err != nil {
		return err
	}

	fmt.Println("Created Airlines table")
	return nil
}
