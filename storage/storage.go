package storage

import "errors"

type Type string

const Memory Type = "memory"
const Hdd Type = "hdd"

var ErrKeyNotExist = errors.New("KEY_NOT_EXIST")
var ErrInvalidStorageType = errors.New("STORAGE_TYPE_INVALID")
