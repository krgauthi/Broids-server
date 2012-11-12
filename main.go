package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Starting broids server")

	nl, err := net.Listen("tcp", "0.0.0.0:9988")
	defer nl.Close()
	if err != nil {
		fmt.Println("main", err)
		return
	}

	m := StartGameManager(nl)
	if m == nil {
		fmt.Println("GameManager:", "failed to start")
		return
	}

	m.NewGame("broids", 1, 160.0, 100.0, "")

	m.Listen()
}
