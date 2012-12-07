package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

type GameManager struct {
	listener net.Listener

	games     map[string]*Game
	gamesLock sync.Mutex
}

func StartGameManager(bind string) *GameManager {
	gm := &GameManager{}

	var err error
	gm.listener, err = net.Listen("tcp", bind)
	if err != nil {
		return nil
	}

	gm.games = make(map[string]*Game)

	go gm.StartGameCleaner()

	return gm
}

func (gm *GameManager) JoinGame(c *Client, name, pass string) {
	g, ok := gm.games[name]
	if !ok {
		c.SendError("game doesn't exist")
		return
	}

	if g.PlayerCount() >= g.limit {
		c.SendError("game is full")
		return
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	out := Frame{Command: FRAME_LOBBY_JOIN}
	temp := JoinOutputData{Id: g.NextId(c), Host: false}

	c.entities = make(map[string]Entity)

	// TODO: Send newplayer to everyone else

	// TODO: Handle host on disconnect
	if g.PlayerCount() == 1 {
		temp.Host = true
	}

	out.Data = temp

	c.encoder.Encode(out)

	fmt.Println(temp)

	g.players[strconv.Itoa(temp.Id)] = c
	c.game = g
}

func (gm *GameManager) NewGame(c *Client, name string, limit int, x, y float32, pass string) {
	g := &Game{}

	g.players = make(map[string]*Client)

	g.name = name
	g.password = pass

	g.limit = limit
	g.height = x
	g.width = y

	g.players["1"] = &Client{Id: 1, entities: make(map[string]Entity)}

	gm.gamesLock.Lock()
	defer gm.gamesLock.Unlock()

	if _, exists := gm.games[g.name]; exists {
		c.SendError("game already exists")
		return
	}

	gm.games[g.name] = g

	if c != nil {
		// TODO: Find a better solution other than running async
		go gm.JoinGame(c, g.name, g.password)
	}
}

func (gm *GameManager) Listen() {
	for {
		conn, err := gm.listener.Accept()
		fmt.Println("Connect")
		if err != nil {
			fmt.Println(err)
		} else {
			go gm.HandleConnect(conn)
		}
	}
}

func (gm *GameManager) ListGames(c *Client) {
	temp := Frame{Command: FRAME_LOBBY_LIST}
	gm.gamesLock.Lock()
	list := make([]ListOutputData, 0)
	for k := range gm.games {
		cur := gm.games[k]
		out := ListOutputData{Name: cur.name, Limit: cur.limit, Current: len(cur.players), Private: cur.Private()}
		list = append(list, out)
	}
	temp.Data = list

	fmt.Println(list)

	c.encoder.Encode(temp)
	gm.gamesLock.Unlock()
}

func (gm *GameManager) StartGameCleaner() {
	in := time.Tick(30 * time.Second)
	for {
		<-in
		gm.gamesLock.Lock()
		for k := range gm.games {
			g := gm.games[k]
			if g.PlayerCount() == 0 {
				fmt.Println("Removing game for inactivity")
				delete(gm.games, k)
			}
		}
		gm.gamesLock.Unlock()
	}
}

func (gm *GameManager) HandleConnect(c net.Conn) {
	cl := &Client{}

	// Save this so we can .close() it later
	cl.conn = c

	// This is where the magic happens
	cl.encoder = json.NewEncoder(c)
	cl.decoder = json.NewDecoder(c)

	go cl.Handle(gm)
}
