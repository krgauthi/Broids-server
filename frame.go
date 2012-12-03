package main

import (
	"encoding/json"
)

// Everything in this file is essentially the spec of what we will be getting and sending

type OutputCommand int

const (
	FRAME_ERROR OutputCommand = -1
	FRAME_SYNC                = 1

	// Delta commands must be == related InputCommands
	FRAME_ENTITY_UPDATE = 2
	FRAME_ENTITY_REMOVE = 3
	FRAME_ENTITY_CREATE = 4
	FRAME_PLAYER_REMOVE = 5
	FRAME_PLAYER_CREATE = 6
	FRAME_HOST_CHANGE   = 7

	// Responses to Lobby Commands
	FRAME_LIST_RESPONSE   = 10
	FRAME_CREATE_RESPONSE = 11
	FRAME_JOIN_RESPONSE   = 12
	FRAME_LEAVE_RESPONSE  = 13
)

type InputCommand int

const (
	COMMAND_EOF   InputCommand = -2
	COMMAND_ERROR              = -1

	// Game Commands
	COMMAND_LEAVE         = 1
	COMMAND_ENTITY_UPDATE = 2
	COMMAND_ENTITY_REMOVE = 3
	COMMAND_ENTITY_CREATE = 4
	COMMAND_PLAYER_REMOVE = 5
	COMMAND_PLAYER_CREATE = 6
	COMMAND_SYNC_REQUEST  = 7

	// Lobby Commands
	COMMAND_LIST   = 10
	COMMAND_CREATE = 11
	COMMAND_JOIN   = 12
)

type EntityType int

const (
	ENTITY_SHIP     EntityType = 1
	ENTITY_ASTEROID            = 2
	ENTITY_BULLET              = 3
)

type SyncFrame struct {
	Command  OutputCommand `json:"c"`
	GameTime int           `json:"t"`
	Data     []Entity      `json:"e"`
}

type DeltaFrame struct {
	Command  OutputCommand `json:"c"`
	GameTime int           `json:"t"`
	Data     Entity        `json:"e"`
}

type InputFrame struct {
	Command InputCommand    `json:"c"`
	Data    json.RawMessage `json:"d"`
}

type JoinInputFrame struct {
	Name     string `json:"n"`
	Password string `json:"p"`
}

type ListOutputFrame struct {
	Command OutputCommand         `json:"c"`
	Data    []ListOutputFrameData `json:"d"`
}

type LeaveOutputFrame struct {
	Command OutputCommand `json:"c"`
}

type ListOutputFrameData struct {
	Name    string `json:"n"`
	Private int    `json:"p"`
	Max     int    `json:"m"`
	Current int    `json:"c"`
}

type ErrorOutputFrame struct {
	Command OutputCommand `json:"c"`
	Code    int           `json:"id"`
	Text    string        `json:"text"`
}

type CreateOutputFrame struct {
	Command OutputCommand `json:"c"`
	Data    string        `json:"d"`
}

type HostOutputFrame struct {
	Command OutputCommand `json:"c"`
	Data    string        `json:"d"`
}

type JoinOutputFrame struct {
	Command OutputCommand       `json:"c"`
	Data    JoinOutputFrameData `json:"d"`
}

type JoinOutputFrameData struct {
	Id     int     `json:"id"`
	Width  float32 `json:"w"`
	Height float32 `json:"h"`
}
