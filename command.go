package main

type CommandId int

const (
	USERNAME CommandId = iota
	JOIN
	ROOM
	MSG
	QUIT
)

type command struct {
	id     CommandId
	client *client
	arg    []string
}
