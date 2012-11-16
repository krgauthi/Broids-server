package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type GameManager struct {
	games    map[string]*Game // A map of all games, keyed by name
	listener net.Listener     // The listener being used
	clients  []*Client        // A list of all the clients

	clientsLock sync.Mutex // Lock when modifying clients
	gamesLock   sync.Mutex // Lock when modifying games
}

// Make sure all games can access the main game manager
var gm *GameManager = nil

// Initialize the game manager
func StartGameManager(nl net.Listener) *GameManager {
	gm = &GameManager{}
	gm.listener = nl
	gm.games = make(map[string]*Game)

	return gm
}

// Listen for new clients and send them off, concurrently to be handled
func (gm *GameManager) Listen() {
	for {
		conn, err := gm.listener.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			go gm.Handle(conn)
		}
	}
}

// Handle new clients
func (gm *GameManager) Handle(c net.Conn) {
	cl := &Client{}
	cl.Entities = make(map[string]*Entity)

	// We save this so we can close() it later
	cl.conn = c

	// These are what we'll be reading from
	cl.decoder = json.NewDecoder(c)
	cl.encoder = json.NewEncoder(c)

	// Add us to the list
	gm.clientsLock.Lock()
	gm.clients = append(gm.clients, cl)
	gm.clientsLock.Unlock()

	go cl.Handle()
}
