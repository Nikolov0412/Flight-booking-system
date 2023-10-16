package main

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type FlightSection struct {
	ID        string `json:"id"`
	SeatClass string `json:"seatClass"`
	NumRows   int    `json:"numRows"`
	NumCols   int    `json:"numCols"`
}

func CreateFlightSection(flightSection FlightSection, svc *dynamodb.DynamoDB) error {
	// Generate a unique ID for the flight section.
	flightSectionID := uuid.New().String()

	if !doesTableExist("FlightSections", svc) {
		if err := createFlightSectionsTable(svc); err != nil {
			fmt.Printf("Error creating FlightSections table: %v\n", err)
		}
	}
	// Create a DynamoDB PutItem input.
	input := &dynamodb.PutItemInput{
		TableName: aws.String("FlightSections"),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(flightSectionID),
			},
			"SeatClass": {
				S: aws.String(flightSection.SeatClass),
			},
			"NumRows": {
				N: aws.String(fmt.Sprintf("%d", flightSection.NumRows)),
			},
			"NumCols": {
				N: aws.String(fmt.Sprintf("%d", flightSection.NumCols)),
			},
		},
	}

	_, err := svc.PutItem(input)
	if err != nil {
		return err
	}

	fmt.Printf("Created Flight Section: ID=%s, SeatClass=%s, NumRows=%d, NumCols=%d\n", flightSectionID, flightSection.SeatClass, flightSection.NumRows, flightSection.NumCols)
	return nil
}

func GetAllFlightSections(svc *dynamodb.DynamoDB) ([]FlightSection, error) {
	// Create a DynamoDB Scan input to scan the FlightSections table.
	input := &dynamodb.ScanInput{
		TableName: aws.String("FlightSections"),
	}

	// Perform the scan operation.
	result, err := svc.Scan(input)
	if err != nil {
		return nil, err
	}

	// Unmarshal the results into a slice of FlightSection.
	var flightSections []FlightSection
	for _, item := range result.Items {
		var flightSection FlightSection
		err := dynamodbattribute.UnmarshalMap(item, &flightSection)
		if err != nil {
			return nil, err
		}

		flightSections = append(flightSections, flightSection)
	}

	return flightSections, nil
}
func GetFlightSectionByID(sectionID string, svc *dynamodb.DynamoDB) (*FlightSection, error) {
	// Create a DynamoDB GetItem input.
	input := &dynamodb.GetItemInput{
		TableName: aws.String("FlightSections"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(sectionID),
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
		return nil, errors.New("Flight Section not found")
	}

	// Parse the retrieved data into a FlightSection struct.
	flightSection := &FlightSection{}
	err = dynamodbattribute.UnmarshalMap(result.Item, flightSection)
	if err != nil {
		return nil, err
	}

	return flightSection, nil
}

func createFlightSectionsTable(svc *dynamodb.DynamoDB) error {
	params := &dynamodb.CreateTableInput{
		TableName: aws.String("FlightSections"),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}
	_, err := svc.CreateTable(params)
	if err != nil {
		return err
	}

	fmt.Println("FlightSections table created successfully")
	return nil
}
