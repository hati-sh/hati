package core

// {MESSAGE_HEADER}CL:{MESSAGE_CONTENT_LENGTH}\n{MESSAGE_EXTRA_SPACE}\n
// {COMMAND}{COMMAND_DELIMITER} <- content
// {COMMAND}{COMMAND_DELIMITER} <- content
// {COMMAND}{COMMAND_DELIMITER} <- content
// {MESSAGE_EOF}

const READ_BLOCK_SIZE = 4096

// VERSION current protocol version eg. v0.1 = {'0', '1'}
var VERSION [2]byte = [2]byte{'0', '1'}

// MESSAGE_HEADER expected header at the beginning of each message
var MESSAGE_HEADER [8]byte = [8]byte{'+', '+', 'h', 'a', 't', 'i', VERSION[0], VERSION[1]}

// MESSAGE_CONTENT_LENGTH message payload content length
var MESSAGE_CONTENT_LENGTH uint64

// MESSAGE_EOF expected EOF - indicates end of message
var MESSAGE_EOF [8]byte = [8]byte{'-', '-', 'h', 'a', 't', 'i', '\n', '\r'}

var MESSAGE_EXTRA_SPACE [8]byte = [8]byte{}
var COMMAND_DELIMITER [2]byte = [2]byte{'+', '\n'}
