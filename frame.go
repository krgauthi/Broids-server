package main

import (
	"encoding/json"
)

type FrameType int

const (
	// Server -> Client communication
	FRAME_ERROR              FrameType = -1
	FRAME_GAME_LEAVE                   = 0
	FRAME_GAME_ENTITY_CREATE           = 1
	FRAME_GAME_ENTITY_MODIFY           = 2
	FRAME_GAME_ENTITY_REMOVE           = 3
	FRAME_GAME_COLLISION               = 4
	FRAME_GAME_PLAYER_CREATE           = 5
	FRAME_GAME_PLAYER_MODIFY           = 6
	FRAME_GAME_PLAYER_REMOVE           = 7
	FRAME_GAME_ROUND_OVER              = 8
	FRAME_GAME_SYNC                    = 9

	// The only game frame without a command equivalent
	FRAME_GAME_HOST_CHANGE = 10

	FRAME_LOBBY_LIST   = 20
	FRAME_LOBBY_CREATE = 21
	FRAME_LOBBY_JOIN   = 22
)

type CommandType int

const (
	// Client -> Server communication
	COMMAND_GAME_LEAVE         CommandType = 0
	COMMAND_GAME_ENTITY_CREATE             = 1
	COMMAND_GAME_ENTITY_MODIFY             = 2
	COMMAND_GAME_ENTITY_REMOVE             = 3
	COMMAND_GAME_COLLISION                 = 4
	COMMAND_GAME_PLAYER_CREATE             = 5
	COMMAND_GAME_PLAYER_MODIFY             = 6
	COMMAND_GAME_PLAYER_REMOVE             = 7
	COMMAND_GAME_ROUND_OVER                = 8

	COMMAND_LOBBY_LIST   = 20
	COMMAND_LOBBY_CREATE = 21
	COMMAND_LOBBY_JOIN   = 22
)

// Server -> Client (Output)
type Frame struct {
	Command FrameType   `json:"c"`
	Data    interface{} `json:"d"`
}

type ListOutputData struct {
	Name    string `json:"n"`
	Limit   int    `json:"l"`
	Current int    `json:"c"`
	Private bool   `json:"p"`
}

type SyncOutputData struct {
	Players  []*Client `json:"p"`
	Entities []*Entity `json:"e"`
}

type CollisionOutputData struct {
	EntityA string `json:"a"`
	EntityB string `json:"b"`
}

type DeltaOutputData interface{}

type JoinOutputData struct {
	Id   int     `json:"i"`
	Host bool    `json:"h"`
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
}

// Client -> Server (Input)
type Command struct {
	Command CommandType     `json:"c"`
	Data    json.RawMessage `json:"d"`
}

type CreateInputData struct {
	Name  string  `json:"n"`
	Limit int     `json:"l"`
	X     float32 `json:"x"`
	Y     float32 `json:"y"`
	Pass  string  `json:"p"`
}

type JoinInputData struct {
	Name string `json:"n"`
	Pass string `json:"p"`
}

type CollisionInputData struct {
	EntityA string `json:"a"`
	EntityB string `json:"b"`
}

type EntityCreateInputData *Entity
type EntityModifyInputData *Entity
type EntityRemoveInputData string

type PlayerCreateInputData *Client
type PlayerModifyInputData *Client
type PlayerRemoveInputData string
