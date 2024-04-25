package core

type Command []byte

var CmdOk []byte = []byte("+OK\n")
var CmdErr []byte = []byte("+ERR\n")

var CmdSet []byte = []byte("SET")
var CmdGet []byte = []byte("GET")
var CmdHas []byte = []byte("HAS")
var CmdDelete []byte = []byte("DELETE")

var CmdCreateRouter []byte = []byte("CREATE_ROUTER")
var CmdCreateQueue []byte = []byte("CREATE_QUEUE")

var CmdPublish []byte = []byte("PUBLISH")

var Commands [][]byte = [][]byte{
	CmdSet, CmdGet, CmdHas, CmdDelete, CmdCreateRouter, CmdCreateQueue, CmdPublish,
}
