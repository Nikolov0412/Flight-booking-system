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

func CreatetSeat(row, col int, isBooked bool, svc *dynamodb.DynamoDB) error {
	// Generate a unique ID for the seat.
	seatID := uuid.New().String()

	// Create a DynamoDB PutItem input.
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Seats"),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(seatID),
			},
			"Row": {
				N: aws.String(fmt.Sprintf("%d", row)),
			},
			"Col": {
				N: aws.String(fmt.Sprintf("%d", col)),
			},
			"IsBooked": {
				BOOL: aws.Bool(isBooked),
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
