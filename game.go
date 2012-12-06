package main

import (
	"strconv"
	"strings"
	"sync"
)

type Game struct {
	name     string
	password string

	limit  int
	height float32
	width  float32

	players map[string]*Client

	// Because we use a single lock,
	// we can use most calls like they're not async,
	// making the whole solution a lot nicer to deal with
	lock sync.Mutex

	lastPlayerId int
}

type EntityType int

const (
	ENTITY_ASTEROID EntityType = 0
	ENTITY_SHIP                = 1
	ENTITY_BULLET              = 2
)

type Entity struct {
	Type EntityType `json:"t"`
	Id   string     `json:"id"`
	X    float32    `json:"x"`
	Y    float32    `json:"y"`
	Xv   float32    `json:"xv"`
	Yv   float32    `json:"yv"`
	A    float32    `json:"a"`
	Av   float32    `json:"av"`
}

func (g *Game) SyncFrame(c *Client) {
	// TODO: This should be in client somewhere
	if c.game == nil {
		c.SendError("no game assigned to player")
		return
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	out := Frame{Command: FRAME_GAME_SYNC}
	temp := SyncOutputData{Entities: nil, Players: nil}
	pout := make([]*Client, 0)
	eout := make([]*Entity, 0)
	for k := range g.players {
		pout = append(pout, g.players[k])
		for ek := range g.players[k].entities {
			eout = append(eout, g.players[k].entities[ek])
		}
	}

	temp.Players = pout
	temp.Entities = eout
	out.Data = temp

	c.encoder.Encode(temp)
}

func (g *Game) CreateEntity(e *Entity) {
	var wg sync.WaitGroup

	g.lock.Lock()
	defer g.lock.Unlock()

	out := Frame{Command: FRAME_GAME_ENTITY_CREATE}
	out.Data = e

	for k := range g.players {
		wg.Add(1)
		go func() {
			p := g.players[k]
			p.encoder.Encode(out)
			wg.Done()
		}()
	}

	wg.Wait()

	// TODO: Error checking
	idParts := strings.SplitN(e.Id, "-", 2)
	c := g.players[idParts[0]]
	c.entities[idParts[1]] = e
}

func (g *Game) ModifyEntity(e *Entity) {
	var wg sync.WaitGroup

	g.lock.Lock()
	defer g.lock.Unlock()

	out := Frame{Command: FRAME_GAME_ENTITY_MODIFY}
	out.Data = e

	for k := range g.players {
		wg.Add(1)
		go func() {
			p := g.players[k]
			p.encoder.Encode(out)
			wg.Done()
		}()
	}

	wg.Wait()

	// TODO: Error checking
	idParts := strings.SplitN(e.Id, "-", 2)
	c := g.players[idParts[0]]
	c.entities[idParts[1]] = e
}

func (g *Game) RemoveEntity(id string) {
	var wg sync.WaitGroup

	g.lock.Lock()
	defer g.lock.Unlock()

	out := Frame{Command: FRAME_GAME_ENTITY_REMOVE}
	out.Data = id

	for k := range g.players {
		wg.Add(1)
		go func() {
			p := g.players[k]
			p.encoder.Encode(out)
			wg.Done()
		}()
	}

	wg.Wait()

	idParts := strings.SplitN(id, "-", 2)
	c := g.players[idParts[0]]
	delete(c.entities, idParts[1])
}

func (g *Game) CreatePlayer(c *Client) {
	var wg sync.WaitGroup

	g.lock.Lock()
	defer g.lock.Unlock()

	out := Frame{Command: FRAME_GAME_PLAYER_CREATE}
	out.Data = c

	for k := range g.players {
		wg.Add(1)
		go func() {
			p := g.players[k]
			p.encoder.Encode(out)
			wg.Done()
		}()
	}

	wg.Wait()
}

func (g *Game) ModifyPlayer(c *Client) {
	var wg sync.WaitGroup

	g.lock.Lock()
	defer g.lock.Unlock()

	out := Frame{Command: FRAME_GAME_PLAYER_MODIFY}
	out.Data = c

	for k := range g.players {
		wg.Add(1)
		go func() {
			p := g.players[k]
			p.encoder.Encode(out)
			wg.Done()
		}()
	}

	wg.Wait()
}

func (g *Game) RemovePlayer(id string) {
	var wg sync.WaitGroup

	g.lock.Lock()
	defer g.lock.Unlock()

	out := Frame{Command: FRAME_GAME_PLAYER_REMOVE}
	out.Data = id

	for k := range g.players {
		wg.Add(1)
		go func() {
			p := g.players[k]
			p.encoder.Encode(out)
			wg.Done()
		}()
	}

	wg.Wait()
}

func (g *Game) RoundOver() {
	// TODO: Implement
}

func (g *Game) Leave(c *Client) {
	g.lock.Lock()
	defer g.lock.Unlock()

	// TODO: Remove player

	var frames sync.WaitGroup
	for ek := range c.entities {
		frames.Add(1)
		go func() {
			var players sync.WaitGroup
			for pk := range g.players {
				players.Add(1)
				go func() {

					frames.Done()
				}()
			}
			frames.Done()
		}()
	}
}

func (g *Game) Collision(a, b string) {
	// TODO: Implement
}

func (g *Game) DeltaFrame(c FrameType, data *DeltaOutputData) {
	var wg sync.WaitGroup

	g.lock.Lock()
	defer g.lock.Unlock()

	out := Frame{Command: c}

	for k := range g.players {
		wg.Add(1)
		go func() {
			p := g.players[k]
			p.encoder.Encode(out)
			wg.Done()
		}()
	}

	wg.Wait()

}

func (g *Game) NextId(c *Client) int {
	ret := g.lastPlayerId
	for {
		ret++

		// 2 is the first available player id - 1 is comp and 0 is temp
		if ret < 2 {
			ret = 2
		}

		if _, ok := g.players[strconv.Itoa(ret)]; !ok {
			c.Id = ret
			g.players[strconv.Itoa(ret)] = c
			break
		}
	}

	g.lastPlayerId = ret

	return ret
}

func (g *Game) Private() bool {
	return len(g.password) != 0
}
