package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type GameManager struct {
	games    map[string]*Game
	listener net.Listener
	clients  []*Client

	clientsLock sync.Mutex
	gamesLock   sync.Mutex
}

var gm *GameManager = nil

func StartGameManager(nl net.Listener) *GameManager {
	gm = &GameManager{}
	gm.listener = nl
	gm.games = make(map[string]*Game)

	return gm
}

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

func (gm *GameManager) Handle(c net.Conn) {
	cl := &Client{}
	cl.Entities = make(map[string]*Entity)
	cl.conn = c
	cl.decoder = json.NewDecoder(c)
	cl.encoder = json.NewEncoder(c)

	gm.clientsLock.Lock()
	defer gm.clientsLock.Unlock()

	gm.clients = append(gm.clients, cl)

	go cl.Handle()
}
