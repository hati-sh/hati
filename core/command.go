package core

type Command []byte

var CmdOk = []byte("+OK\n")
var CmdErr = []byte("+ERR\n")

var CmdSet = []byte("SET")
var CmdGet = []byte("GET")
var CmdHas = []byte("HAS")
var CmdCount = []byte("COUNT")
var CmdDelete = []byte("DELETE")
var CmdFlushAll = []byte("FLUSHALL")

var CmdCreateRouter = []byte("CREATE_ROUTER")
var CmdCreateQueue = []byte("CREATE_QUEUE")

var CmdPublish = []byte("PUBLISH")
var CmdPurge = []byte("PURGE")
var CmdGetRouter = []byte("GET_ROUTER")
var CmdGetQueue = []byte("GET_QUEUE")
var CmdListRouter = []byte("LIST_ROUTER")
var CmdListQueue = []byte("LIST_QUEUE")

var Commands = [][]byte{
	CmdSet, CmdGet, CmdHas, CmdDelete, CmdCreateRouter, CmdCreateQueue, CmdPublish, CmdFlushAll, CmdPurge,
	CmdGetRouter, CmdGetQueue, CmdListRouter, CmdListQueue,
}
