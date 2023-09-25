package main

type FlightSection struct {
	ID        int      `json:"id"`
	SeatClass string   `json:"seatClass"`
	NumRows   int      `json:"numRows"`
	NumCols   int      `json:"numCols"`
	Seats     [][]Seat `json:"seats"`
}
