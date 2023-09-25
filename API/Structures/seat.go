package main

type Seat struct {
	ID       int  `json:"id"`
	Row      int  `json:"row"`
	Col      int  `json:"col"`
	IsBooked bool `json:"isBooked"`
}
