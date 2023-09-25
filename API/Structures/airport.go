package main

import (
	"errors"
	"regexp"
	"strings"
)

type Airport struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
}

func ValidateAirportName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) != 3 || !isAlphabetic(name) {
		return errors.New("Airport name must be exactly 3 alphabetic characters")
	}
	return nil
	// TODO: Implement validation for unique name when the database is implemented.
}

// isAlphabetic checks if a string contains only alphabetic characters.
func isAlphabetic(s string) bool {
	return regexp.MustCompile("^[a-zA-Z]+$").MatchString(s)
}
