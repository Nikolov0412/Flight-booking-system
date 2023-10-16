package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type Flight struct {
	ID                 string        `json:"id"`
	FlightNumber       string        `json:"flightNumber"`
	FlightSectionID    []string      `json:flightSectionID`
	OriginAirport      string        `json:"originAirport"`
	DestinationAirport string        `json:"destinationAirport"`
	DepartureDate      time.Time     `json:"departureDate"`
	FlightTime         time.Duration `json:"flightTime"`
	ETA                string        `json:"eta"`
}

func CreateFlight(flight Flight, svc *dynamodb.DynamoDB) error {
	// Validate the flight data as needed.
	if err := validateFlightData(flight); err != nil {
		return err
	}
	if !doesTableExist("Flights", svc) {
		if err := createFlightsTable(svc); err != nil {
			fmt.Printf("Error creating Flights table: %v\n", err)
		}
	}
	flightID := uuid.New().String()
	// Convert the list of FlightSectionID strings to a list of DynamoDB attribute values.
	flightSectionIDs := make([]*string, len(flight.FlightSectionID))
	for i, id := range flight.FlightSectionID {
		flightSectionIDs[i] = aws.String(id)
	}

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
			"FlightSectionID": {
				SS: flightSectionIDs,
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
				S: aws.String(CalculateETA(flight)),
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

func GetAllFlights(svc *dynamodb.DynamoDB) ([]Flight, error) {
	// Create a DynamoDB Scan input to scan the "Flights" table.
	input := &dynamodb.ScanInput{
		TableName: aws.String("Flights"),
	}

	// Perform the scan operation.
	result, err := svc.Scan(input)
	if err != nil {
		return nil, err
	}

	// Unmarshal the results into a slice of Flight.
	var flights []Flight
	for _, item := range result.Items {
		var flight Flight
		err := dynamodbattribute.UnmarshalMap(item, &flight)
		if err != nil {
			return nil, err
		}

		flights = append(flights, flight)
	}

	return flights, nil
}
func GetFlightsByOriginAirport(originAirport string, svc *dynamodb.DynamoDB) ([]Flight, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Flights"),
		IndexName:              aws.String("originAirport"), // Use the index name you defined
		KeyConditionExpression: aws.String("OriginAirport = :originAirport"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":originAirport": {
				S: aws.String(originAirport),
			},
		},
	}

	result, err := svc.Query(input)
	if err != nil {
		return nil, err
	}

	var flights []Flight
	for _, item := range result.Items {
		var flight Flight
		if err := dynamodbattribute.UnmarshalMap(item, &flight); err != nil {
			return nil, err
		}
		flights = append(flights, flight)
	}

	return flights, nil
}
func GetFlightsByDestinationAirport(destinationAirport string, svc *dynamodb.DynamoDB) ([]Flight, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Flights"),
		IndexName:              aws.String("destinationAirport"), // Use the index name you defined
		KeyConditionExpression: aws.String("DestinationAirport = :DestinationAirport"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":DestinationAirport": {
				S: aws.String(destinationAirport),
			},
		},
	}

	result, err := svc.Query(input)
	if err != nil {
		return nil, err
	}

	var flights []Flight
	for _, item := range result.Items {
		var flight Flight
		if err := dynamodbattribute.UnmarshalMap(item, &flight); err != nil {
			return nil, err
		}
		flights = append(flights, flight)
	}

	return flights, nil
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
		{flight.FlightSectionID, "FlightSectionID is required"},
		{flight.OriginAirport, "OriginAirport is required"},
		{flight.DestinationAirport, "DestinationAirport is required"},
		{flight.DepartureDate, "DepartureDate is required and must be a valid date"},
		{flight.FlightTime, "FlightTime must be greater than 0"},
	}

	for _, rule := range validationRules {
		if !isNotEmpty(rule.field) {
			return errors.New(rule.message)
		}
	}

	return nil
}

// CalculateETA calculates the Estimated Time of Arrival (ETA) and returns it as a formatted string.
func CalculateETA(flight Flight) string {
	departureTime := flight.DepartureDate
	arrivalTime := departureTime.Add(flight.FlightTime)

	// Format the ETA as "15:28" (HH:mm).
	etaString := arrivalTime.Format("15:28")

	return etaString
}

func createFlightsTable(svc *dynamodb.DynamoDB) error {
	// Define the parameters for creating the "Flights" table.
	params := &dynamodb.CreateTableInput{
		TableName: aws.String("Flights"),
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
				IndexName: aws.String("originAiport"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("OriginAirport"),
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
			{
				IndexName: aws.String("destinationAirport"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("DestinationAirport"),
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

	// Create the "Flights" table.
	_, err := svc.CreateTable(params)
	if err != nil {
		return err
	}

	fmt.Println("Flights table created successfully")
	return nil
}
