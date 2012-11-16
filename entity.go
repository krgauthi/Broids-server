package main

type Entity struct {
	Id              string     `json:"id"` // Entity Id
	Type            EntityType `json:"t"`  // Entity type
	XPos            float32    `json:"x"`  // X position of the entity
	YPos            float32    `json:"y"`  // Y position of the entity
	Ang             float32    `json:"a"`  // Angle of the entity
	AngularVelocity float32    `json:"av"` // Angular velocity of the entity
	LinXVelocity    float32    `json:"xv"` // Linear velocity in the X direction
	LinYVelocity    float32    `json:"yv"` // Linear velocity in the Y direction
}
