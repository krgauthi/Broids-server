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
}

// This will store all the games
var games map[string]*Game

func init() {
	games = make(map[string]*Game)
}

func GameNew(nl net.Listener, name string, count int) *Game {
	game := &Game{MaxPlayers: count, Name: name}
	game.players = make(map[string]*Player)
	game.lock = &sync.RWMutex{}
	game.done = make(chan bool)
	game.playerAccept = make(chan net.Conn)

	return game
}

func (g *Game) SyncFrames() {
	timer := time.Tick(1 * time.Second)
	for {
		select {
		case <-timer:
			fmt.Println("Tick!")
			g.GameTime++
			s := Frame{GameTime: g.GameTime}
			for pid := range g.players {
				p := g.players[pid]
				for eid := range p.Entities {
					// We dereference so it copies the struct
					e := *p.Entities[eid]
					e.Id = p.Name + "-" + e.Id
					s.Frames = append(s.Frames, e)
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

func (p *Player) jsonChan() chan Frame {
	j := make(chan Frame)
	go func() {
		dec := json.NewDecoder(p.conn)
		for {
			var m Frame
			if err := dec.Decode(&m); err == io.EOF {
				// TODO: End of IO
				j <- Frame{Type: -1}
				break
			} else if err != nil {
				j <- Frame{Type: -1}
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
				for i := range p.Entities {
					p.Entities[i] = nil
					delete(p.Entities, i)
				}

				g.players[p.Name] = nil
				delete(g.players, p.Name)

				fmt.Println("Player disconnected")

				// TODO: Send removal frame

				return
			}
			d, err := json.Marshal(f)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// TODO: Squash
			for p := range g.players {
				g.players[p].conn.Write(d)
				g.players[p].conn.Write([]byte("\n"))
			}
		}
	}
}

func (g *Game) Stop() {
	// TODO: Shut down go routines
}
