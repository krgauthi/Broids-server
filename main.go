package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Starting broids server")

	// Start the server, listening on all interfaces, on port 9988
	nl, err := net.Listen("tcp", "0.0.0.0:9988")
	defer nl.Close()
	if err != nil {
		fmt.Println("main", err)
		return
	}

	// Start the game manager
	m := StartGameManager(nl)
	if m == nil {
		fmt.Println("GameManager:", "failed to start")
		return
	}

	// Create a new game, mostly for testing purposes
	m.NewGame("broids", 1, 160.0, 100.0, "")

	// Handle new connections
	m.Listen()
}
