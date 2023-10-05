package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type FlightSection struct {
	ID        string   `json:"id"`
	SeatClass string   `json:"seatClass"`
	NumRows   int      `json:"numRows"`
	NumCols   int      `json:"numCols"`
	Seats     [][]Seat `json:"seats"`
}

func CreateFlightSection(flightSection FlightSection, svc *dynamodb.DynamoDB) error {
	// Generate a unique ID for the flight section.
	flightSectionID := uuid.New().String()

	// Flatten the Seats array into a single list of seats.
	var seats []Seat
	for _, row := range flightSection.Seats {
		seats = append(seats, row...)
	}

	// Convert the seats to a list of DynamoDB attribute values.
	seatAVList := make([]*dynamodb.AttributeValue, len(seats))
	for i, seat := range seats {
		seatAV, err := dynamodbattribute.MarshalMap(seat)
		if err != nil {
			return err
		}
		seatAVList[i] = &dynamodb.AttributeValue{M: seatAV}
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
			"Seats": {
				L: seatAVList, // Store the flattened seats as a list of maps.
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

		// Unmarshal the Seats attribute (list of maps) into a slice of Seat.
		for _, av := range item["Seats"].L {
			var seat Seat
			err := dynamodbattribute.UnmarshalMap(av.M, &seat)
			if err != nil {
				return nil, err
			}
			flightSection.Seats = append(flightSection.Seats, []Seat{seat}) // Append as a slice of slices
		}

		flightSections = append(flightSections, flightSection)
	}

	return flightSections, nil
}
