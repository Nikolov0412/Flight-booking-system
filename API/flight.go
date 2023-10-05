package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
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
	flightID := uuid.New().String()
	// Create a new flight item in DynamoDB.
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String("Flights"),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(flightID),
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

	fmt.Printf("Created Flight: ID=%s, FlightNumber=%s\n", flightID, flight.FlightNumber)
	return nil
}

func isNotEmpty(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return v != ""
	case time.Time:
		return !v.IsZero()
	case int:
		return v != 0
	// Add more cases for other types as needed
	default:
		return true // Assume not empty for unsupported types
	}
}

func validateFlightData(flight Flight) error {
	var validationRules = []struct {
		field   interface{}
		message string
	}{
		{flight.FlightNumber, "FlightNumber is required"},
		{flight.OriginAirport, "OriginAirport is required"},
		{flight.DestinationAirport, "DestinationAirport is required"},
		{flight.DepartureDate, "DepartureDate is required and must be a valid date"},
		{flight.FlightTime, "FlightTime must be greater than 0"},
		{flight.ETA, "ETA is required"},
		{flight.PlaneID, "PlaneID is required"},
	}

	for _, rule := range validationRules {
		if !isNotEmpty(rule.field) {
			return errors.New(rule.message)
		}
	}

	return nil
}

// CalculateETA calculates the Estimated Time of Arrival (ETA) and returns it as a formatted string.
func (flight Flight) CalculateETA() string {
	departureTime := flight.DepartureDate
	arrivalTime := departureTime.Add(flight.FlightTime)

	// Format the ETA as "15:28" (HH:mm).
	etaString := arrivalTime.Format("15:28")

	return etaString
}
