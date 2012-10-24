package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Starting broids server")

	nl, err := net.Listen("tcp", "0.0.0.0:9988")
	defer nl.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	game := GameNew(nl, "broids", 5)
	if game == nil {
		os.Exit(1)
	}
	game.Start()

	// TODO: Cleanly kill the server
	for {
		conn, err := nl.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			go AcceptConnection(conn)
		}
	}
}

func AcceptConnection(c net.Conn) {
	cf := &ConnectFrame{}
	dec := json.NewDecoder(c)
	dec.Decode(cf)

	g, ok := games[cf.Game]
	if !ok {
		fmt.Println("Game", cf.Game, "does not exist")
		// TODO: Send crap saying game doesn't exist
	} else {
		fmt.Println("Connection for", g.Name)
		// TODO: Send crap to say accepted
		g.playerAccept <- c
	}
}
