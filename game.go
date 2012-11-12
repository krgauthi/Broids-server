package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Game struct {
	Name       string
	Width      float32
	Height     float32
	MaxPlayers int
	GameTime   int // Not really used

	lastId     int
	password   string
	deltaStore map[string]*DeltaFrame
	players    map[int]*Client

	joinLock sync.Mutex
	sendLock sync.Mutex

	stop      chan bool
	syncStop  chan bool
	deltaStop chan bool
}

type GameError struct {
	Code int
	Text string
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

func (g *Game) Private() bool {
	return len(g.password) != 0
}

func (g *Game) Start() {
	go g.SyncFrames()
	go g.DeltaFrames()
}

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
		}
	}
}

func (g *Game) SyncFrames() {
	syncTimer := time.Tick(1 * time.Second)

	for {
		select {
		case <-syncTimer:
			g.sendLock.Lock()
			fmt.Println(g.Name+":", "Sync!")
			g.GameTime++
			s := SyncFrame{Command: FRAME_SYNC, GameTime: g.GameTime}
			for pid := range g.players {
				p := g.players[pid]
				for eid := range p.Entities {
					// We dereference so it copies the struct
					d := *p.Entities[eid]
					d.Id = d.Id
					s.Data = append(s.Data, d)
					fmt.Println(d)
				}
			}

			data, err := json.Marshal(s)
			if err != nil {
				continue
			}

			for p := range g.players {
				g.players[p].conn.Write(data)
				g.players[p].conn.Write([]byte("\n"))
			}

			// Clear the delta
			for i := range g.deltaStore {
				g.deltaStore[i] = nil
				delete(g.deltaStore, i)
			}
			g.sendLock.Unlock()
		}
	}
}

func (p *Client) Join(name string) *GameError {
	g, ok := gm.games[name]
	if !ok {
		return &GameError{Code: 1, Text: "game doesn't exist"}
	}

	if len(g.players) >= g.MaxPlayers {
		return &GameError{Code: 2, Text: "game full"}
	}

	g.joinLock.Lock()
	defer g.joinLock.Unlock()

	nextId := g.lastId

	// NOTE: This will shit itself if we have more players than int >= 0 can handle
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

	p.Status = STATUS_GAME
	p.game = g

	return nil
}
