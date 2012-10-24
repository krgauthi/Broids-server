package main

type Entity struct {
	Id        string  `json:"id"`
	Type      int     `json:"t"`
	XPos      float64 `json:"x"`
	YPos      float64 `json:"y"`
	Direction float64 `json:"d"`
	Velocity  float64 `json:"v"`
}

const (
	TYPE_SHIP = iota
	TYPE_ASTEROID
)
