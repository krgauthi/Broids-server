package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

type Game struct {
	Name         string
	MaxPlayers   int
	LastId       int
	GameTime     int
	MapX         float32
	Mapy         float32
	players      map[string]*Player
	lock         *sync.RWMutex
	playerAccept chan net.Conn
	done         chan bool
	delta        map[string]*ActionData
}

// This will store all the games
var games map[string]*Game

func init() {
	games = make(map[string]*Game)
}

func GameNew(nl net.Listener, name string, count int) *Game {
	game := &Game{MaxPlayers: count, Name: name}
	game.players = make(map[string]*Player)
	game.delta = make(map[string]*ActionData)
	game.lock = &sync.RWMutex{}
	game.done = make(chan bool)
	game.playerAccept = make(chan net.Conn)

	return game
}

func (g *Game) SyncFrames() {
	syncTimer := time.Tick(1 * time.Second)
	deltaTimer := time.Tick(100 * time.Millisecond)

	for {
		select {
		case <-syncTimer:
			g.lock.Lock()
			fmt.Println("Tick!")
			g.GameTime++
			s := Frame{Type: FRAME_SYNC, GameTime: g.GameTime}
			for pid := range g.players {
				p := g.players[pid]
				for eid := range p.Entities {
					// We dereference so it copies the struct
					e := *p.Entities[eid]
					e.Id = p.Name + "-" + e.Id
					d := ActionData{Type: ACTION_MOVE, Entity: e}
					s.Data = append(s.Data, d)
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
			for i := range g.delta {
				g.delta[i] = nil
				delete(g.delta, i)
			}
			g.lock.Unlock()
		case <-deltaTimer:
			g.lock.Lock()
			if len(g.delta) > 0 {
				fmt.Println("Delta!")
				s := Frame{Type: FRAME_DELTA, GameTime: g.GameTime}

				for e := range g.delta {
					s.Data = append(s.Data, *g.delta[e])
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
				for i := range g.delta {
					g.delta[i] = nil
					delete(g.delta, i)
				}
			}
			g.lock.Unlock()
		}
	}
}

func (g *Game) Start() {
	go g.SyncFrames()
	go g.AcceptPlayers()
	/*if games == nil {
		fmt.Println("Well, shit")
	}*/
	games[g.Name] = g
}

func (g *Game) AcceptPlayers() {
	for {
		select {
		case p := <-g.playerAccept:
			// TODO: Fix issue if overflow in g.LastId
			temp := &Player{
				Name:     strconv.Itoa(g.LastId),
				conn:     p,
				Entities: make(map[string]*Entity),
			}
			g.LastId++
			g.players[temp.Name] = temp
			go g.HandleClient(temp)
		}
	}
}

func (p *Player) jsonChan() chan ActionData {
	j := make(chan ActionData)
	go func() {
		dec := json.NewDecoder(p.conn)
		for {
			var m ActionData
			if err := dec.Decode(&m); err == io.EOF {
				// TODO: End of IO
				j <- ActionData{Type: -1}
				break
			} else if err != nil {
				j <- ActionData{Type: -1}
				return
			}

			j <- m
		}
	}()
	return j
}

func (g *Game) HandleClient(p *Player) {
	defer p.conn.Close()

	fmt.Printf("%s: cid: %s\n", g.Name, p.Name)

	s := &Entity{Type: TYPE_SHIP, Id: strconv.Itoa(p.LastId)}
	// TODO: Fix issue if overflow in p.LastId
	p.Entities[s.Id] = s
	p.LastId++

	g.lock.Lock()
	e := *p.Entities[s.Id]
	e.Id = p.Name + "-" + e.Id
	g.delta[e.Id] = &ActionData{Type: ACTION_CREATE, Entity: e}
	g.lock.Unlock()

	j := p.jsonChan()

	// TODO: Rewrite this
	for {
		select {
		case <-p.quit:
			// TODO: Notify client
			break
		case f := <-j:
			if f.Type == -1 {
				// Err reading - remove this player
				// Clear out entities - this is for GC stuff
				g.lock.Lock()
				for i := range p.Entities {
					e := *p.Entities[i]
					e.Id = p.Name + "-" + e.Id
					g.delta[p.Entities[i].Id] = &ActionData{Type: ACTION_DESTROY, Entity: e}
					p.Entities[i] = nil
					delete(p.Entities, i)
				}
				g.lock.Unlock()

				g.players[p.Name] = nil
				delete(g.players, p.Name)

				fmt.Println("Player disconnected")

				return
			}

			g.lock.Lock()
			g.delta[f.Entity.Id] = &f
			g.lock.Unlock()
		}
	}
}

func (g *Game) Stop() {
	// TODO: Shut down go routines
}
