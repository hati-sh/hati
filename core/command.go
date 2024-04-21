package core

type Cmd []byte

var CMD_SET Cmd = []byte("SET")
var CMD_GET Cmd = []byte("GET")

type Command struct {
	Cmd   Cmd
	Key   string
	Value string
}
