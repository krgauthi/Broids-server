package main

import (
	"fmt"
	//"strconv"
)

func protocolVersion() int {
	return 2
}

func main() {
	fmt.Println("Starting Broids Server")

	// Start server, listening on all interfaces on port 9988
	m := StartGameManager("0.0.0.0:9988")
	if m == nil {
		fmt.Println("main:", "failed to start GameManager")
		return
	}

	// Testing stuff
	/*for i := 0; i < 13; i++ {
		m.NewGame(nil, "broids"+strconv.Itoa(i), 5, 100.0, 160.0, "")
	}*/
	m.NewGame(nil, "broids", 5, 100.0, 160.0, "")

	// Handle new connections
	m.Listen()
}
