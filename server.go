package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type server struct {
	rooms    map[string]*rooms
	commands chan command
}

func Newserver() *server {

	return &server{

		rooms:    make(map[string]*rooms),
		commands: make(chan command),
	}
}
func (s *server) run() {

	for cmd := range s.commands {
		switch cmd.id {
		case JOIN:
			s.joinRoom(cmd.client, cmd.arg)
		case USERNAME:
			s.username(cmd.client, cmd.arg)
		case MSG:
			s.msg(cmd.client, cmd.arg)

		case QUIT:
			s.quit(cmd.client, cmd.arg)

		}
	}
}
func (s *server) newClient(conn net.Conn) *client {

	log.Printf("new clinet has connected :%s", conn.RemoteAddr().String())
	return &client{
		conn:     conn,
		username: "user",
		commands: s.commands,
	}

}
func (s *server) username(c *client, args []string) {
	if len(args) < 2 {
		c.msg("username is required")
		return
	}
	c.username = args[1]
	c.msg(fmt.Sprintln("New user %s", c.username))

}
func (s *server) joinRoom(c *client, args []string) {
	if len(args) < 2 {
		c.msg("room name is required. usage: /join ROOM_NAME")
		return
	}

	roomName := args[1]

	r, ok := s.rooms[roomName]

	if !ok {
		r = &rooms{
			name:   roomName,
			member: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.member[c.conn.RemoteAddr()] = c

	r.len = r.len + 1
	if r.len > 5 {
		c.msg("full")
		return
	}
	go func() {
		for {
		
			time.Sleep(10 * time.Second)
		    c.msg("10 second")
		}
	}()
	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.username))

}

func (s *server) msg(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("first enter the room first"))
		return
	}
	c.room.broadcast(c, c.username+":+"+strings.Join(args[1:len(args)], " "))
}
func (s *server) quit(c *client, args []string) {
	log.Println("client has disconnected : %s" + strings.Join(args[1:len(args)], " "))
	s.quitCurrentRoom(c)

	c.msg("see you soon")
	c.conn.Close()
}
func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.member, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s haas left the room ", c.username))
	}
}
