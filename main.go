package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting broids server")

	// Start the game manager, listening on all interfaces on port 9988
	m := StartGameManager("0.0.0.0:9988")
	if m == nil {
		fmt.Println("main:", "failed to start GameManager")
		return
	}

	// Create a new game, mostly for testing purposes
	m.NewGame("broids", 1, 160.0, 100.0, "")

	// Handle new connections
	m.Listen()
}
