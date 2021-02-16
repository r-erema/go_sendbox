package hash_table

import (
	"bytes"
	"encoding/binary"
	"go_sendbox/learning/other/DBs/DBil/files_storage"
	"go_sendbox/learning/other/DBs/DBil/generic"
	"go_sendbox/learning/other/DBs/DBil/kv_protocol"
	"io"
)

type hashTable struct {
	indexes      map[string]map[string]int
	filesStorage *files_storage.DataFilesStorage
}

func NewHashTable(dataFiles *files_storage.DataFilesStorage) *hashTable {
	return &hashTable{filesStorage: dataFiles}
}

func (ht *hashTable) Save(key, value string) error {
	kv := kv_protocol.KV{}
	copy(kv.Key[:], key)
	copy(kv.Value[:], value)
	lastFile := ht.filesStorage.LastFile()
	return binary.Write(lastFile, binary.LittleEndian, kv)
}

func (ht *hashTable) Get(key string) (value *string) {
	if ht.indexes != nil {
		for fileName, index := range ht.indexes {
			if offset, ok := index[key]; ok {
				file := ht.filesStorage.GetFileByName(fileName)
				_, _ = file.Seek(int64(offset), io.SeekStart)
				kv := kv_protocol.KV{}
				_ = binary.Read(file, binary.LittleEndian, &kv)
				value := string(bytes.Trim(kv.Value[:], "\x00"))
				return &value
			}
		}
	}

	return generic.FindInFiles(ht.filesStorage.Files(), key)
}

func (ht *hashTable) Index() error {
	kv := kv_protocol.KV{}
	ht.indexes = map[string]map[string]int{}

	for _, file := range ht.filesStorage.Files() {
		currentOffset := 0
		_, _ = file.Seek(0, io.SeekStart)
		ht.indexes[file.Name()] = make(map[string]int)
		for {
			err := binary.Read(file, binary.LittleEndian, &kv)
			if err == io.EOF {
				break
			}
			key := string(bytes.Trim(kv.Key[:], "\x00"))
			ht.indexes[file.Name()][key] = currentOffset
			currentOffset += kv_protocol.KVLengthBytes
		}
	}

	return nil
}

func (ht *hashTable) Clean() error {
	ht.indexes = nil
	return nil
}
