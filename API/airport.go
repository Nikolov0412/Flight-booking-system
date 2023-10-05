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
	ID   string `json:"id"`
	Code string `json:"code"`
}

func ValidateAirportCode(code string) error {
	code = strings.TrimSpace(code)
	if len(code) != 3 || !isAlphabetic(code) {
		return errors.New("Airport code must be exactly 3 alphabetic characters")
	}
	return nil
}

func CreateAirport(airport Airport, svc *dynamodb.DynamoDB) error {
	// Validate the airport code.
	if err := ValidateAirportCode(airport.Code); err != nil {
		return err
	}
	if !doesTableExist("Airports", svc) {
		if err := createAirportsTable(svc); err != nil {
			fmt.Printf("Error creating Airports table: %v\n", err)
		}
	}
	// Check if the airport code is already in use.
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("Airports"),
		IndexName: aws.String("CodeIndex"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":code": {
				S: aws.String(airport.Code),
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
	airportID := uuid.New().String()
	// Create a new airport in DynamoDB with ID and Code.
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String("Airports"),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(airportID),
			},
			"Code": {
				S: aws.String(airport.Code),
			},
		},
	}
	_, err = svc.PutItem(putInput)
	if err != nil {
		return err
	}

	fmt.Printf("Created Airport: ID=%s, Code=%s\n", airportID, airport.Code)
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

func GetAllAirports(svc *dynamodb.DynamoDB) ([]*Airport, error) {
	// Create a DynamoDB Scan input to retrieve all items from the Airports table.
	input := &dynamodb.ScanInput{
		TableName: aws.String("Airports"),
	}

	// Perform the scan operation.
	result, err := svc.Scan(input)
	if err != nil {
		return nil, err
	}

	// Initialize a slice to hold the retrieved airports.
	airports := []*Airport{}

	// Iterate through the scan results and parse each item into an Airport struct.
	for _, item := range result.Items {
		airport := &Airport{}
		if err := dynamodbattribute.UnmarshalMap(item, airport); err != nil {
			return nil, err
		}
		airports = append(airports, airport)
	}

	return airports, nil
}

func createAirportsTable(svc *dynamodb.DynamoDB) error {
	// Define the parameters for creating the "Airports" table.
	params := &dynamodb.CreateTableInput{
		TableName: aws.String("Airports"),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"), // Primary key
			},
			{
				AttributeName: aws.String("Code"), // Secondary key
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
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("CodeIndex"), // Name of the GSI
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("Code"),
						KeyType:       aws.String("HASH"), // GSI key
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"), // Include all attributes in the index
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	// Create the "Airports" table.
	_, err := svc.CreateTable(params)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created Airports table")
	return nil
}
