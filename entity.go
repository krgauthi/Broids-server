package main

type Entity struct {
	Id              string     `json:"id"`
	Type            EntityType `json:"t"`
	XPos            float32    `json:"x"`
	YPos            float32    `json:"y"`
	Ang             float32    `json:"a"`
	AngularVelocity float32    `json:"av"`
	LinXVelocity    float32    `json:"xv"`
	LinYVelocity    float32    `json:"yv"`
}
