package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn `json:"-"`

	game     *Game             `json:"-"`
	Id       int               `json:"i"`
	Name     string            `json:"n"`
	Score    int               `json:"s"`
	Color    string            `json:"c"`
	Host     bool              `json:"h"`
	entities map[string]Entity `json:"-"`

	encoder *json.Encoder `json:"-"`
	decoder *json.Decoder `json:"-"`
}

func (c *Client) SendError(err string) {
	frame := Frame{Command: FRAME_ERROR}
	frame.Data = err
	c.encoder.Encode(frame)
}

func (c *Client) Disconnect() {
	if c.game != nil {
		c.game.Leave(c)
	}

	c.conn.Close()
}

func (c *Client) Handle(gm *GameManager) {
	reallyExit := false

	for {
		var command Command
		err := c.decoder.Decode(&command)
		if err != nil {
			c.Disconnect()
			break
		}

		switch command.Command {
		case COMMAND_LOBBY_CREATE:
			var in CreateInputData
			json.Unmarshal(command.Data, &in)
			gm.NewGame(c, in.Name, in.Limit, in.X, in.Y, in.Pass)
		case COMMAND_LOBBY_JOIN:
			var in JoinInputData
			json.Unmarshal(command.Data, &in)
			gm.JoinGame(c, in.Name, in.Pass)
		case COMMAND_LOBBY_LIST:
			gm.ListGames(c)
		}

		if c.game != nil {
			for {
				err = c.decoder.Decode(&command)
				if err != nil {
					fmt.Println("BYE")
					c.game.RemovePlayer(string(c.Id))
					c.Disconnect()
					reallyExit = true
					break
				}

				switch command.Command {
				case COMMAND_GAME_LEAVE:
					fmt.Println("LEAVE")
					c.game.Leave(c)
					break
				case COMMAND_GAME_ENTITY_CREATE:
					fmt.Println("ENTITY CREATE")
					var in EntityCreateInputData
					json.Unmarshal(command.Data, &in)
					fmt.Println(in)
					c.game.CreateEntity(Entity(in))
				case COMMAND_GAME_ENTITY_MODIFY:
					fmt.Println("ENTITY MODIFY")
					var in EntityModifyInputData
					json.Unmarshal(command.Data, &in)
					c.game.ModifyEntity(Entity(in))
				case COMMAND_GAME_ENTITY_REMOVE:
					fmt.Println("ENTITY REMOVE")
					var in EntityRemoveInputData
					json.Unmarshal(command.Data, &in)
					fmt.Println(in)
					c.game.RemoveEntity(string(in))
				case COMMAND_GAME_COLLISION:
					fmt.Println("BOOM")
					var in CollisionInputData
					json.Unmarshal(command.Data, &in)
					c.game.Collision(in.EntityA, in.EntityB)
				case COMMAND_GAME_PLAYER_CREATE:
					var in PlayerCreateInputData
					json.Unmarshal(command.Data, &in)
					c.game.CreatePlayer(in)
				case COMMAND_GAME_PLAYER_MODIFY:
					var in PlayerModifyInputData
					json.Unmarshal(command.Data, &in)
					c.game.ModifyPlayer(in)
				case COMMAND_GAME_PLAYER_REMOVE:
					var in PlayerRemoveInputData
					json.Unmarshal(command.Data, &in)
					c.game.RemovePlayer(string(in))
				case COMMAND_GAME_ROUND_OVER:
					c.game.RoundOver()
				}
			}
		}

		if reallyExit {
			break
		}
	}
}
