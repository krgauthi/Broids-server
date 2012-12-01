package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
)

type ClientStatus int

const (
	STATUS_LOBBY ClientStatus = 1
	STATUS_GAME               = 2
)

type Client struct {
	Name     string
	Id       int
	Status   ClientStatus
	Entities map[string]*Entity

	game    *Game
	conn    net.Conn
	decoder *json.Decoder
	encoder *json.Encoder

	quit chan bool
}

// This is magic
func (c *Client) jsonChan() chan InputFrame {
	j := make(chan InputFrame)
	go func() {
		for {
			var m InputFrame
			if err := c.decoder.Decode(&m); err == io.EOF {
				// TODO: End of IO
				j <- InputFrame{Command: COMMAND_EOF}
				return
			} else if err != nil {
				j <- InputFrame{Command: COMMAND_ERROR}
				return
			}

			j <- m
		}
	}()
	return j
}

func (c *Client) Leave() {
	// Err reading - remove this player
	// Clear out entities - this is for GC stuff
	if c.game != nil {
		c.game.sendLock.Lock()
		for i := range c.Entities {
			e := *c.Entities[i]
			e.Id = c.Name + "-" + e.Id
			c.game.deltaStore[c.Entities[i].Id] = &DeltaFrame{Command: COMMAND_ENTITY_REMOVE}
			// TODO: Copy entity
			c.Entities[i] = nil
			delete(c.Entities, i)
		}

		c.game.players[c.Id] = nil
		delete(c.game.players, c.Id)

		c.game.sendLock.Unlock()
		c.game = nil
		c.Status = STATUS_LOBBY

		fmt.Println("Player left game")
	}

	return
}

// Send an error to the client, based on a GameError
func (c *Client) Error(err *GameError) {
	c.encoder.Encode(ErrorOutputFrame{Command: FRAME_ERROR, Text: err.Text, Code: err.Code})
}

func (c *Client) Sync() {
	c.game.sendLock.Lock()

	fmt.Println(c.game.Name+":", "Sync!")
	c.game.GameTime++
	s := SyncFrame{Command: FRAME_SYNC, GameTime: c.game.GameTime}

	// Create a list of all the entities
	for pid := range c.game.players {
		p := c.game.players[pid]
		for eid := range p.Entities {
			// We dereference so it copies the struct
			d := *p.Entities[eid]
			d.Id = d.Id
			s.Data = append(s.Data, d)
			fmt.Println(d)
		}
	}

	// Prepare the data for sending
	data, err := json.Marshal(s)
	if err != nil {
		// Keep calm and error on
		return
	}

	// Send the sync frame to this user
	c.conn.Write(data)
	c.conn.Write([]byte("\n"))

	// Clear the delta
	for i := range c.game.deltaStore {
		c.game.deltaStore[i] = nil
		delete(c.game.deltaStore, i)
	}
	c.game.sendLock.Unlock()
}

func (c *Client) Handle() {
	// TODO: Rewrite this

	j := c.jsonChan()

	for {
		select {
		case <-c.quit:
			// TODO: Notify client
			break
		case f := <-j:
			if c.Status == STATUS_GAME {
				switch f.Command {
				case COMMAND_EOF, COMMAND_ERROR:
					// If in a game, on error, leave the game and close the connection
					c.Leave()
					c.conn.Close()
				case COMMAND_ENTITY_CREATE, COMMAND_ENTITY_REMOVE, COMMAND_ENTITY_UPDATE:
					// Anything that would generate a delta frame
					temp := DeltaFrame{}
					temp.Command = OutputCommand(f.Command)
					json.Unmarshal(f.Data, &temp.Data)
					c.game.sendLock.Lock()

					// Update or remove?
					switch f.Command {
					case COMMAND_ENTITY_CREATE, COMMAND_ENTITY_UPDATE:
						fmt.Println("Create/Update")
						ent := Entity(temp.Data)
						c.Entities[temp.Data.Id] = &ent
					case COMMAND_ENTITY_REMOVE:
						fmt.Println("Delete")
						delete(c.Entities, temp.Data.Id)
					}
					c.game.deltaStore[temp.Data.Id] = &temp

					c.game.sendLock.Unlock()
				case COMMAND_REQUEST_SYNC:
					c.Sync()
				case COMMAND_LEAVE:
					c.Leave()
				}
			} else {
				switch f.Command {
				case COMMAND_EOF:
					c.conn.Close()
				case COMMAND_JOIN:
					temp := JoinInputFrame{}
					json.Unmarshal(f.Data, &temp)

					// Join the game
					err := c.Join(temp.Name)
					if err != nil {
						c.Error(err)
					}
				case COMMAND_LIST:
					temp := ListOutputFrame{Command: COMMAND_LIST}
					temp.Data = make([]ListOutputFrameData, 0)
					for k := range gm.games {
						game := gm.games[k]
						var priv int
						if game.Private() {
							priv = 1
						}

						temp.Data = append(temp.Data, ListOutputFrameData{
							Name:    game.Name,
							Current: len(game.players),
							Max:     game.MaxPlayers,
							Private: priv},
						)
					}
					c.encoder.Encode(temp)
				}
			}
		}
	}
}
