package main

type Frame struct {
	Type     FrameType    `json:"t"`
	GameTime int          `json:"gt"`
	Data     []ActionData `json:"d"`
}

type ConnectFrame struct {
	Game string `json:"g"`
}

type ActionData struct {
	Type   ActionType `json:"t"`
	Entity Entity     `json:"e"`
}

type ActionType int

const (
	ACTION_CREATE  ActionType = 1
	ACTION_DESTROY            = 2
	ACTION_MOVE               = 3
)

type FrameType int

const (
	FRAME_DELTA FrameType = 1
	FRAME_SYNC            = 2
)
