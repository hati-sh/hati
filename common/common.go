package common

import (
	"bytes"
	"encoding/gob"
)

func EncodeToBytes(input any) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(input)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
