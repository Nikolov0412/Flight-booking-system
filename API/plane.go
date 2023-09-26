package main

type Plane struct {
	ID             int             `json:"id"`
	Name           string          `json:"name"`
	AirlineID      int             `json:"airlineId"`
	FlightSections []FlightSection `json:"flightSections"`
}
