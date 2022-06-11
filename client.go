package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	username string
	room     *rooms
	commands chan<- command
}

func (c *client) readinput() {
	for {

		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}
		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {

		case "/username":
			c.commands <- command{
				id:     USERNAME,
				client: c,
				arg:    args,
			}

		case "/msg":
			c.commands <- command{
				id:     MSG,
				client: c,
				arg:    args,
			}
		case "/join":
			c.commands <- command{
				id:   JOIN,
				client: c,
				arg:    args,
			}
		case "/quit":
			c.commands <- command{
				id:     QUIT,
				client: c,
				arg:    args,
			}

		default:
			c.err(fmt.Errorf("check commands%s", cmd))
		}
	}
}
func (c *client) err(err error) {
	c.conn.Write([]byte("ERR:" + err.Error() + "\n"))
}
func (c *client) msg(msg string) {
	c.conn.Write([]byte(">" + msg + "\n"))
}
