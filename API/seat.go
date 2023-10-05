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
	Row      int  `json:"Row"`
	Col      int  `json:"Col"`
	IsBooked bool `json:"IsBooked"`
}

func CreatetSeat(seat Seat, svc *dynamodb.DynamoDB) error {
	seatID := uuid.New().String()

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
