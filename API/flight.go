package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Flight struct {
	ID                 string        `json:"id"`
	FlightNumber       string        `json:"flightNumber"`
	OriginAirport      string        `json:"originAirport"`
	DestinationAirport string        `json:"destinationAirport"`
	DepartureDate      time.Time     `json:"departureDate"`
	FlightTime         time.Duration `json:"flightTime"`
	ETA                string        `json:"eta"`
	PlaneID            string        `json:"planeID"` // Reference to the Plane in the Planes table
}

func CreateFlight(flight Flight, svc *dynamodb.DynamoDB) error {
	// Validate the flight data as needed.
	if err := validateFlightData(flight); err != nil {
		return err
	}

	// Create a new flight item in DynamoDB.
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String("Flights"),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(flight.ID),
			},
			"FlightNumber": {
				S: aws.String(flight.FlightNumber),
			},
			"OriginAirport": {
				S: aws.String(flight.OriginAirport),
			},
			"DestinationAirport": {
				S: aws.String(flight.DestinationAirport),
			},
			"DepartureDate": {
				S: aws.String(flight.DepartureDate.Format(time.RFC3339)), // Store as a string
			},
			"FlightTime": {
				N: aws.String(fmt.Sprintf("%d", flight.FlightTime.Milliseconds())), // Store as milliseconds
			},
			"ETA": {
				S: aws.String(flight.ETA),
			},
			"PlaneID": {
				S: aws.String(flight.PlaneID),
			},
		},
	}

	_, err := svc.PutItem(putInput)
	if err != nil {
		return err
	}

	fmt.Printf("Created Flight: ID=%s, FlightNumber=%s\n", flight.ID, flight.FlightNumber)
	return nil
}

func validateFlightData(flight Flight) error {
	// TODO: Implement any data validation rules here.
	// For example, check that FlightNumber, OriginAirport, and DestinationAirport are not empty, etc.
	if flight.FlightNumber == "" {
		return errors.New("FlightNumber is required")
	}
	return nil
}
