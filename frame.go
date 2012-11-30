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
	FRAME_DELTA_UPDATE = 2
	FRAME_DELTA_REMOVE = 3
	FRAME_DELTA_CREATE = 4

	// Responses to Lobby Commands
	FRAME_LIST_RESPONSE   = 10
	FRAME_CREATE_RESPONSE = 11
	FRAME_JOIN_RESPONSE   = 12
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
	COMMAND_REQUEST_SYNC  = 5

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

type ListOutputFrameData struct {
	Name    string `json:"n"`
	Private int    `json:"p"`
	Max     int    `json:"m"`
	Current int    `json:"c"`
}

type ErrorOutputFrame struct {
	Command OutputCommand `json:"t"`
	Code    int           `json:"id"`
	Text    string        `json:"text"`
}

type JoinOutputFrame struct {
	Command OutputCommand `json:"c"`
	Data    int           `json:"d"`
}
