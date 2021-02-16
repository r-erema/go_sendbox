package generic

import (
	"bytes"
	"encoding/binary"
	"go_sendbox/learning/other/DBs/DBil/kv_protocol"
	"io"
	"os"
)

func FindInFiles(files []*os.File, key string) (value *string) {
	kv := kv_protocol.KV{}
	var found []string
	for _, file := range files {
		_, _ = file.Seek(0, io.SeekStart)
		for {
			err := binary.Read(file, binary.LittleEndian, &kv)
			if err == io.EOF {
				break
			}
			keyFromFile := string(bytes.Trim(kv.Key[:], "\x00"))
			if keyFromFile == key {
				value := string(bytes.Trim(kv.Value[:], "\x00"))
				found = append(found, value)
			}
		}
	}

	if len(found) > 0 {
		return &found[len(found)-1]
	}

	return nil
}
