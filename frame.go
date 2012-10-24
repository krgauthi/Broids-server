package main

type Frame struct {
	Type     int      `json:"t"`
	GameTime int      `json:"gt"`
	Frames   []Entity `json:"f"`
}

type ConnectFrame struct {
	Game string `json:"gid"`
}
