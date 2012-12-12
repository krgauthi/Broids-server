package main

import (
	"fmt"
	//"os"
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
	Type  EntityType `json:"t"`
	Id    string     `json:"id"`
	X     float32    `json:"x"`
	Y     float32    `json:"y"`
	Xv    float32    `json:"xv"`
	Yv    float32    `json:"yv"`
	A     float32    `json:"a"`
	Av    float32    `json:"av"`
	Extra int        `json:"e"`
}

// TODO: Send HOST change when needed

func (g *Game) SyncFrame(c *Client) {
	out := &Frame{Command: FRAME_GAME_SYNC}
	temp := SyncOutputData{Entities: nil, Players: nil}
	pout := make([]*Client, 0)
	eout := make([]*Entity, 0)
	for k := range g.players {
		pout = append(pout, g.players[k])
		for ek := range g.players[k].entities {
			temp := g.players[k].entities[ek]
			eout = append(eout, &temp)
		}
	}

	temp.Players = pout
	temp.Entities = eout
	out.Data = temp

	c.encoder.Encode(out)
}

func (g *Game) SendFrame(f *Frame) {
	var wg sync.WaitGroup
	for k := range g.players {
		p := g.players[k]
		if p.encoder != nil {
			wg.Add(1)
			go func() {
				p.encoder.Encode(f)
				wg.Done()
			}()
		}
	}

	wg.Wait()
}

func (g *Game) CreateEntity(e Entity) {
	g.lock.Lock()
	defer g.lock.Unlock()

	fmt.Println(g.name, "new entity", e.Id)

	out := &Frame{Command: FRAME_GAME_ENTITY_CREATE}
	out.Data = e

	g.SendFrame(out)

	// TODO: Error checking
	idParts := strings.SplitN(e.Id, "-", 2)
	c := g.players[idParts[0]]
	c.entities[idParts[1]] = e
}

func (g *Game) ModifyEntity(e Entity) {
	g.lock.Lock()
	defer g.lock.Unlock()

	fmt.Println(g.name, "modify entity", e.Id)

	out := &Frame{Command: FRAME_GAME_ENTITY_MODIFY}
	out.Data = e

	g.SendFrame(out)

	// TODO: Error checking
	idParts := strings.SplitN(e.Id, "-", 2)
	c := g.players[idParts[0]]
	c.entities[idParts[1]] = e
}

func (g *Game) RemoveEntity(id string) {
	g.lock.Lock()
	defer g.lock.Unlock()

	fmt.Println(g.name, "delete entity", id)

	out := &Frame{Command: FRAME_GAME_ENTITY_REMOVE}
	out.Data = id

	g.SendFrame(out)

	idParts := strings.SplitN(id, "-", 2)
	c := g.players[idParts[0]]
	delete(c.entities, idParts[1])
}

func (g *Game) CreatePlayer(c *Client) {
	g.lock.Lock()
	defer g.lock.Unlock()

	fmt.Println(g.name, "new player", c.Id)

	out := &Frame{Command: FRAME_GAME_PLAYER_CREATE}
	out.Data = c

	g.SendFrame(out)
}

func (g *Game) ModifyPlayer(c *Client) {
	g.lock.Lock()
	defer g.lock.Unlock()

	fmt.Println(g.name, "modify player", c.Id)

	out := &Frame{Command: FRAME_GAME_PLAYER_MODIFY}
	out.Data = c

	g.SendFrame(out)
}

func (g *Game) RemovePlayer(id string) {
	g.lock.Lock()
	defer g.lock.Unlock()

	fmt.Println(g.name, "remove player", id)

	out := &Frame{Command: FRAME_GAME_PLAYER_REMOVE}
	out.Data = id

	g.SendFrame(out)
}

func (g *Game) RoundOver() {
	g.lock.Lock()
	defer g.lock.Unlock()

	fmt.Println(g.name, "round over")

	out := &Frame{Command: FRAME_GAME_ROUND_OVER}

	g.SendFrame(out)
}

func (g *Game) Leave(c *Client) {
	g.lock.Lock()
	defer g.lock.Unlock()

	fmt.Println(g.name, "player leaves:", c.Name)

	if _, ok := g.players[strconv.Itoa(c.Id)]; ok {
		delete(g.players, strconv.Itoa(c.Id))
	}

	var frames sync.WaitGroup
	for ek := range c.entities {
		frames.Add(1)
		out := &Frame{Command: FRAME_GAME_ENTITY_REMOVE}
		out.Data = c.entities[ek].Id
		go func() {
			g.SendFrame(out)
			frames.Done()
		}()
	}
	frames.Wait()

	leave := &Frame{Command: FRAME_GAME_LEAVE}
	c.encoder.Encode(leave)
}

func (g *Game) Collision(a string, ap int, b string, bp int) {
	g.lock.Lock()
	defer g.lock.Unlock()

	fmt.Println(g.name, "collision:", a, b)

	out := &Frame{Command: FRAME_GAME_COLLISION}
	d := CollisionOutputData{EntityA: a, APoints: ap, EntityB: b, BPoints: bp}
	out.Data = d

	g.SendFrame(out)
}

func (g *Game) DeltaFrame(c FrameType, data *DeltaOutputData) {
	g.lock.Lock()
	defer g.lock.Unlock()

	out := &Frame{Command: c}

	g.SendFrame(out)
}

func (g *Game) PlayerCount() int {
	return len(g.players) - 1
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
