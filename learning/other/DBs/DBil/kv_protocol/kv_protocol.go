package kv_protocol

const (
	keyLengthBytes   = 30
	valueLengthBytes = 100
	KVLengthBytes    = keyLengthBytes + valueLengthBytes
)

type KV struct {
	Key   [keyLengthBytes]byte
	Value [valueLengthBytes]byte
}
