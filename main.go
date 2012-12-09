package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting Broids Server")

	// Start server, listening on all interfaces on port 9988
	m := StartGameManager("0.0.0.0:9988")
	if m == nil {
		fmt.Println("main:", "failed to start GameManager")
		return
	}

	// Testing stuff
	m.NewGame(nil, "broids", 5, 100.0, 160.0, "")

	// Handle new connections
	m.Listen()
}
