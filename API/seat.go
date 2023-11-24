package main

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type Seat struct {
	ID              string `json:"id"`
	Row             int    `json:"Row"`
	Col             int    `json:"Col"`
	IsBooked        bool   `json:"IsBooked"`
	FlightSectionID string `json:FlightSectionId`
	FlightNumber    string `json:FlightNumber`
}

func CreateSeat(seat Seat, svc *dynamodb.DynamoDB) error {
	seatID := uuid.New().String()

	if !doesTableExist("Seats", svc) {
		if err := createSeatsTable(svc); err != nil {
			fmt.Printf("Error creating Seats table: %v\n", err)
		}
	}

	if err := validateFlightNumber(seat.FlightNumber, svc); err != nil {
		return err
	}
	if err := validateFlightSectionID(seat.FlightSectionID, svc); err != nil {
		return err
	}
	if err := validateRowColInFlightSection(seat, svc); err != nil {
		return err
	}

	// Create a DynamoDB PutItem input.
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Seats"),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(seatID),
			},
			"Row": {
				N: aws.String(fmt.Sprintf("%d", seat.Row)),
			},
			"Col": {
				N: aws.String(fmt.Sprintf("%d", seat.Col)),
			},
			"FlightNumber": {
				S: aws.String(seat.FlightNumber),
			},
			"FlightSectionID": {
				S: aws.String(seat.FlightSectionID),
			},
			"IsBooked": {
				BOOL: aws.Bool(seat.IsBooked),
			},
		},
	}

	// Insert the item into DynamoDB.
	_, err := svc.PutItem(input)
	if err != nil {
		return err
	}

	// No errors occurred, return nil.
	return nil
}
func GetSeatsByFlightNumber(FlightNumber string, svc *dynamodb.DynamoDB) ([]*Seat, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String("Seats"),
		IndexName:              aws.String("FlightNumberIndex"), // Use the GSI name
		KeyConditionExpression: aws.String("#FlightNumber = :FlightNumber"),
		ExpressionAttributeNames: map[string]*string{
			"#FlightNumber": aws.String("FlightNumber"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":FlightNumber": {
				S: aws.String(FlightNumber),
			},
		},
	}

	result, err := svc.Query(queryInput)
	if err != nil {
		return nil, err
	}

	seats := []*Seat{}

	for _, item := range result.Items {
		seat := &Seat{}
		if err := dynamodbattribute.UnmarshalMap(item, seat); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}

	return seats, nil
}
func GetSeatsByFlightSectionID(FlightSectionID string, svc *dynamodb.DynamoDB) ([]*Seat, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String("Seats"),
		IndexName:              aws.String("FlightSectionIndex"),
		KeyConditionExpression: aws.String("#FlightSectionID = :FlightSectionID"),
		ExpressionAttributeNames: map[string]*string{
			"#FlightSectionID": aws.String("FlightSectionID"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":FlightSectionID": {
				S: aws.String(FlightSectionID),
			},
		},
	}

	result, err := svc.Query(queryInput)
	if err != nil {
		return nil, err
	}

	seats := []*Seat{}

	for _, item := range result.Items {
		seat := &Seat{}
		if err := dynamodbattribute.UnmarshalMap(item, seat); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}

	return seats, nil
}

func GetSeatByID(seatID string, svc *dynamodb.DynamoDB) (*Seat, error) {
	// Create a DynamoDB GetItem input.
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Seats"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(seatID),
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
		return nil, errors.New("Seat not found")
	}

	// Parse the retrieved data into a Seat struct.
	seat := &Seat{}
	err = dynamodbattribute.UnmarshalMap(result.Item, seat)
	if err != nil {
		return nil, err
	}

	return seat, nil
}
func GetAllSeats(svc *dynamodb.DynamoDB) ([]Seat, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Seats"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		return nil, err
	}

	seats := []Seat{}

	for _, item := range result.Items {
		seat := Seat{}
		if err := dynamodbattribute.UnmarshalMap(item, &seat); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}

	return seats, nil
}

func UpdateSeatIsBooked(seatID, flightSectionID string, isBooked bool, svc *dynamodb.DynamoDB) error {
	// Create a DynamoDB UpdateItem input to update the IsBooked property.
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Seats"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(seatID),
			},
			"FlightSectionID": {
				S: aws.String(flightSectionID),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":isBooked": {
				BOOL: aws.Bool(isBooked),
			},
		},
		UpdateExpression: aws.String("SET IsBooked = :isBooked"),
		ReturnValues:     aws.String("UPDATED_NEW"),
	}

	// Update the IsBooked property of the seat in DynamoDB.
	result, err := svc.UpdateItem(input)
	if err != nil {
		fmt.Printf("Error updating seat %s IsBooked: %v\n", seatID, err)
		return err
	}

	// Check the result for debugging purposes (optional).
	fmt.Printf("UpdateItem result: %v\n", result)

	fmt.Printf("Updated seat %s IsBooked to %v\n", seatID, isBooked)
	return nil
}

func createSeatsTable(svc *dynamodb.DynamoDB) error {
	// Define the parameters for creating the "Seats" table.
	params := &dynamodb.CreateTableInput{
		TableName: aws.String("Seats"),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("FlightSectionID"),
				KeyType:       aws.String("RANGE"),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("FlightSectionID"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("FlightNumber"),
				AttributeType: aws.String("S"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("FlightSectionIndex"), // Name of the GSI
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("FlightSectionID"),
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
			{
				IndexName: aws.String("FlightNumberIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("FlightNumber"),
						KeyType:       aws.String("HASH"),
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
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	// Create the "Seats" table.
	_, err := svc.CreateTable(params)
	if err != nil {
		return err
	}

	fmt.Println("Seats table created successfully")
	return nil
}
func validateFlightNumber(flightNumber string, svc *dynamodb.DynamoDB) error {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String("Flights"),
		IndexName:              aws.String("FlightNumberIndex"),
		KeyConditionExpression: aws.String("#fn = :fn"),
		ExpressionAttributeNames: map[string]*string{
			"#fn": aws.String("FlightNumber"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":fn": {
				S: aws.String(flightNumber),
			},
		},
	}

	// Perform the query.
	result, err := svc.Query(queryInput)
	if err != nil {
		return err
	}

	// If the count is 0, FlightNumber does not exist.
	if *result.Count == 0 {
		return errors.New("FlightNumber does not exist")
	}

	return nil
}

func validateFlightSectionID(flightSectionID string, svc *dynamodb.DynamoDB) error {
	// Query the FlightSections table to check if the FlightSectionID exists.
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String("FlightSections"),
		KeyConditionExpression: aws.String("#id = :id"),
		ExpressionAttributeNames: map[string]*string{
			"#id": aws.String("ID"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(flightSectionID),
			},
		},
	}

	queryResult, err := svc.Query(queryInput)
	if err != nil {
		return err
	}

	if *queryResult.Count == 0 {
		return errors.New("FlightSectionID does not exist")
	}

	return nil
}

func validateRowColInFlightSection(seat Seat, svc *dynamodb.DynamoDB) error {
	// Retrieve the FlightSection details using FlightSectionID from the seat
	flightSection, err := GetFlightSectionByID(seat.FlightSectionID, svc)
	if err != nil {
		return err
	}

	// Check if the provided Row and Col are within the valid range
	if seat.Row < 1 || seat.Row > flightSection.NumRows || seat.Col < 1 || seat.Col > flightSection.NumCols {
		return errors.New("Row or Col is out of range for the FlightSection")
	}

	return nil
}
