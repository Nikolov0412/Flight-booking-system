package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
