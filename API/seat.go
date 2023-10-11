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
		IndexName:              aws.String("FlightSectionIDIndex"), // Use the GSI name
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

func UpdateSeatIsBooked(seatID string, isBooked bool, svc *dynamodb.DynamoDB) error {
	// Create a DynamoDB UpdateItem input to update the IsBooked property.
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Seats"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(seatID),
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
	_, err := svc.UpdateItem(input)
	if err != nil {
		return err
	}

	fmt.Printf("Updated seat %s IsBooked to %v\n", seatID, isBooked)
	return nil
}

func CreateSeatMatrix(numRows, numCols int) [][]Seat {
	seatMatrix := make([][]Seat, numRows)

	for i := 0; i < numRows; i++ {
		seatMatrix[i] = make([]Seat, numCols)
		for j := 0; j < numCols; j++ {
			// Initialize each seat in the matrix
			seatMatrix[i][j] = Seat{
				Row:      i + 1, // Rows and columns are usually 1-based
				Col:      j + 1,
				IsBooked: false,
			}
		}
	}

	return seatMatrix
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
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
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
