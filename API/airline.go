package main

import (
	"errors"
	"strings"
)

type Airline struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func ValidateAirlineName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) >= 6 {
		return errors.New("Airline name must have a length less than 6 characters")
	}
	return nil

	// TODO: Implement the rest validation for unique name after the database is implemented.
}
