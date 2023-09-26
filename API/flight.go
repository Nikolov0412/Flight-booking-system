package main

import "time"

type Flight struct {
	ID                 int           `json:"id"`
	FlightNumber       string        `json:"flightNumber"`
	OriginAirport      string        `json:"originAirport"`
	DestinationAirport string        `json:"destinationAirport"`
	DepartureDate      time.Time     `json:"departureDate"`
	FlightTime         time.Duration `json:"flightTime"`
	ETA                string        `json:"eta"`
	Plane              Plane         `json:"plane"`
}
