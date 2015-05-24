package toolkit

import (
	"gopkg.in/vmihailenco/msgpack.v2"
)

func ObjectToByte(i interface{}) ([]byte, error) {
	return msgpack.Marshal(i)
}

func ByteToObject(data []byte, i interface{}) error {
	return msgpack.Unmarshal(data, i)
}
