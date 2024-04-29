package core

type Command []byte

var CmdOk = []byte("+OK\n")
var CmdErr = []byte("+ERR\n")

var CmdSet = []byte("SET")
var CmdGet = []byte("GET")
var CmdHas = []byte("HAS")
var CmdDelete = []byte("DELETE")
var CmdFlushAll = []byte("FLUSHALL")

var CmdCreateRouter = []byte("CREATE_ROUTER")
var CmdCreateQueue = []byte("CREATE_QUEUE")

var CmdPublish = []byte("PUBLISH")

var Commands = [][]byte{
	CmdSet, CmdGet, CmdHas, CmdDelete, CmdCreateRouter, CmdCreateQueue, CmdPublish, CmdFlushAll,
}
