package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

type Plane struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	AirlineID      int             `json:"airlineId"`
	FlightSections []FlightSection `json:"flightSections"`
}

func CreatePlane(name string, airlineID int, flightSections []FlightSection, svc *dynamodb.DynamoDB) error {
	// Generate a unique ID for the plane.
	planeID := uuid.New().String()

	// Create a DynamoDB PutItem input for the Plane.
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Planes"),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(planeID),
			},
			"Name": {
				S: aws.String(name),
			},
			"AirlineID": {
				N: aws.String(fmt.Sprintf("%d", airlineID)),
			},
			"FlightSections": {
				L: []*dynamodb.AttributeValue{},
			},
		},
	}

	// Add flight sections to the item.
	for _, section := range flightSections {
		// Create a DynamoDB Map to represent the FlightSection.
		flightSectionMap := make(map[string]*dynamodb.AttributeValue)
		flightSectionMap["ID"] = &dynamodb.AttributeValue{S: aws.String(uuid.New().String())}
		flightSectionMap["SeatClass"] = &dynamodb.AttributeValue{S: aws.String(section.SeatClass)}
		flightSectionMap["NumRows"] = &dynamodb.AttributeValue{N: aws.String(fmt.Sprintf("%d", section.NumRows))}
		flightSectionMap["NumCols"] = &dynamodb.AttributeValue{N: aws.String(fmt.Sprintf("%d", section.NumCols))}
		flightSectionMap["Seats"] = &dynamodb.AttributeValue{L: []*dynamodb.AttributeValue{}}

		// Append the FlightSection map to the FlightSections list.
		input.Item["FlightSections"].L = append(input.Item["FlightSections"].L, &dynamodb.AttributeValue{M: flightSectionMap})
	}

	// Store the Plane item in DynamoDB.
	_, err := svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
