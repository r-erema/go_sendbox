package mem_table

import (
	"bytes"
	"encoding/binary"
	"go_sendbox/learning/other/DBs/DBil/kv_protocol"
	"io"
	"os"
	"reflect"
	"sort"
)

type KVPairs []map[string]string

type MemTable struct {
	KVs     KVPairs
	logFile *os.File
}

func NewMemTable(logFile *os.File) *MemTable {
	return &MemTable{KVs: []map[string]string{}, logFile: logFile}
}

func (mt *MemTable) Save(key, value string) error {
	kvPair := make(map[string]string, 1)
	kvPair[key] = value
	mt.KVs = append(mt.KVs, kvPair)
	mt.sort()
	kv := &kv_protocol.KV{}
	copy(kv.Key[:], key)
	copy(kv.Value[:], value)
	return binary.Write(mt.logFile, binary.LittleEndian, kv)
}

func (mt *MemTable) RestoreValuesFromLogFile() error {
	kv := kv_protocol.KV{}
	_, err := mt.logFile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	for {
		err = binary.Read(mt.logFile, binary.LittleEndian, &kv)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		kvPair := make(map[string]string, 1)
		kvPair[string(bytes.Trim(kv.Key[:], "\x00"))] = string(bytes.Trim(kv.Value[:], "\x00"))
		mt.KVs = append(mt.KVs, kvPair)
	}
	return nil
}

func (mt *MemTable) Find(key string) *string {
	for _, kv := range mt.KVs {
		if value, ok := kv[key]; ok {
			return &value
		}
	}
	return nil
}

func (mt *MemTable) sort() {
	var keys []string
	for _, kv := range mt.KVs {
		keys = append(keys, reflect.ValueOf(kv).MapKeys()[0].String())
	}
	sort.Strings(keys)

	result := make(KVPairs, len(mt.KVs))

	for i, k := range keys {
		kvPair := make(map[string]string, 1)
		kvPair[k] = *mt.Find(k)
		result[i] = kvPair
	}

	mt.KVs = result
}
