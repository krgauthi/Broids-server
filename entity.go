package main

type Entity struct {
	Id        string     `json:"id"`
	Type      EntityType `json:"t"`
	XPos      float64    `json:"x"`
	YPos      float64    `json:"y"`
	Direction float64    `json:"d"`
	Velocity  float64    `json:"v"`
}

type EntityType int

const (
	TYPE_SHIP EntityType = iota
	TYPE_ASTEROID
)
