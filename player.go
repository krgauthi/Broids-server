package main

import (
	"net"
)

type Player struct {
	conn     net.Conn
	Name     string
	LastId   int
	Entities map[string]*Entity
	Ship     *Entity
	quit     chan bool
}
