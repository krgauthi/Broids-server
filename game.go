package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Game struct {
	Name       string  // Name of the game
	Width      float32 // Width of the map
	Height     float32 // Height of the map
	MaxPlayers int     // Max number of players for this game
	GameTime   int     // Not really used

	lastId     int                    // Last connected id
	password   string                 // Password for the game, MD5 hashed
	deltaStore map[string]*DeltaFrame // Map of all DeltaFrames to be sent out, keyed by Id
	players    map[int]*Client        // Map of all players, keyed by Id. We don't use an array because this could be really sparse.

	joinLock sync.Mutex // Lock when player joins
	sendLock sync.Mutex // Lock when frames are sent

	stop      chan bool // Channel that will recieve data when quitting
	syncStop  chan bool // Channel for stopping sync frames. Probably not needed any more.
	deltaStop chan bool // Channel for stopping delta frames
}

// This defines an error that we can send using json
type GameError struct {
	Code int    `json:"c"` // Error code
	Text string `json:"t"` // Error description
}

func (e *GameError) Error() string {
	return e.Text
}

// NOTE: This is actually a gamemanager function
func (gm *GameManager) NewGame(name string, max int, x, y float32, password string) error {
	gm.gamesLock.Lock()
	defer gm.gamesLock.Unlock()

	if _, ok := gm.games[name]; ok {
		return errors.New("game already exists")
	}

	g := Game{}
	g.Name = name
	g.MaxPlayers = max
	g.deltaStore = make(map[string]*DeltaFrame)
	g.players = make(map[int]*Client)

	g.stop = make(chan bool)
	g.syncStop = make(chan bool)
	g.deltaStop = make(chan bool)

	gm.games[g.Name] = &g
	go g.Start()

	return nil
}

// It's a private game if the password is empty
func (g *Game) Private() bool {
	return len(g.password) != 0
}

// Starting a game, we need to start sending sync frames and delta frames
func (g *Game) Start() {
	go g.DeltaFrames()
}

// This command handles sending delta frames
func (g *Game) DeltaFrames() {
	deltaTimer := time.Tick(100 * time.Millisecond)
	for {
		select {
		case <-deltaTimer:
			g.sendLock.Lock()
			if len(g.deltaStore) > 0 {
				fmt.Println(g.Name+":", "Delta!")

				frames := make([][]byte, 0)

				for e := range g.deltaStore {
					data, err := json.Marshal(g.deltaStore[e])
					if err != nil {
						continue
					}
					frames = append(frames, data)
				}

				for p := range g.players {
					for f := range frames {
						g.players[p].conn.Write(frames[f])
						g.players[p].conn.Write([]byte("\n"))
					}
				}

				// Clear the delta
				for i := range g.deltaStore {
					g.deltaStore[i] = nil
					delete(g.deltaStore, i)
				}
			}
			g.sendLock.Unlock()
		case <-g.deltaStop:
			break
		}
	}
}

func (p *Client) Join(name string) *GameError {
	// If the game doesn't exist, send an error
	g, ok := gm.games[name]
	if !ok {
		return &GameError{Code: 1, Text: "game doesn't exist"}
	}

	// If the game is full, send an error
	if len(g.players) >= g.MaxPlayers {
		return &GameError{Code: 2, Text: "game full"}
	}

	g.joinLock.Lock()

	// Start on the lastId inserted
	nextId := g.lastId

	// This increments nextId until we find one that isn't used
	// NOTE: This will shit itself if we have more players than int >= 0 can handle
	// It'll busy loop until someone disconnects.
	for {
		if _, ok := g.players[nextId]; !ok {
			break
		}

		nextId++
		if nextId < 0 {
			nextId = 0
		}
	}

	g.players[nextId] = p
	p.Id = nextId
	g.lastId = nextId + 1

	// The player is now in a game
	p.Status = STATUS_GAME
	p.game = g

	g.joinLock.Unlock()

	return nil
}
