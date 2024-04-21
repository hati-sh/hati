package core

type Cmd string

const CMD_SET Cmd = "SET"
const CMD_GET Cmd = "GET"

type Command struct {
	Cmd   Cmd
	Key   string
	Value string
}
