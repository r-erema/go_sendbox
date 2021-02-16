package lsm_tree

import (
	"bytes"
	"encoding/binary"
	"go_sendbox/learning/other/DBs/DBil/files_storage"
	"go_sendbox/learning/other/DBs/DBil/generic"
	"go_sendbox/learning/other/DBs/DBil/kv_protocol"
	"go_sendbox/learning/other/DBs/DBil/mem_table"
	"io"
	"os"
	"reflect"
)

type bytesRange struct {
	minByte, maxByte int
}

type LSMTree struct {
	dataFiles       *files_storage.DataFilesStorage
	memTable        *mem_table.MemTable
	memTableMaxSize int
	indexes         map[string]map[string]bytesRange
}

const memTableLogFilePath = "./data/mt.data.bin"

func NewLSMTree(dataFiles *files_storage.DataFilesStorage, memTableMaxSize int) *LSMTree {
	strategy := &LSMTree{
		dataFiles:       dataFiles,
		memTable:        mem_table.NewMemTable(prepareDataFile()),
		memTableMaxSize: memTableMaxSize,
	}
	_ = strategy.memTable.RestoreValuesFromLogFile()
	return strategy
}

func prepareDataFile() (memTableLogFile *os.File) {
	memTableLogFile, _ = os.OpenFile(memTableLogFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	return
}

func (lt *LSMTree) Save(key, value string) error {

	err := lt.memTable.Save(key, value)
	if err != nil {
		return err
	}

	if len(lt.memTable.KVs) >= lt.memTableMaxSize {
		return lt.pushMemTableToFile()
	}

	return nil
}

func (lt *LSMTree) pushMemTableToFile() error {
	file := lt.dataFiles.LastFile()
	kvs := make([]kv_protocol.KV, len(lt.memTable.KVs))
	i := 0
	for _, kv := range lt.memTable.KVs {
		kvToBin := kv_protocol.KV{}
		k := reflect.ValueOf(kv).MapKeys()[0].String()
		copy(kvToBin.Key[:], k)
		copy(kvToBin.Value[:], kv[k])
		kvs[i] = kvToBin
		i++
	}
	err := binary.Write(file, binary.LittleEndian, kvs)
	if err != nil {
		return err
	}
	lt.memTable = mem_table.NewMemTable(prepareDataFile())
	return nil
}

func (lt *LSMTree) Get(key string) *string {

	value := lt.memTable.Find(key)
	if value != nil {
		return value
	}

	firstLetter := string(key[0])
	if lt.indexes != nil {
		for fileName, index := range lt.indexes {
			file := lt.dataFiles.GetFileByName(fileName)
			_, _ = file.Seek(int64(index[firstLetter].minByte), io.SeekStart)
			kv := kv_protocol.KV{}
			for {
				err := binary.Read(file, binary.LittleEndian, &kv)
				if err == io.EOF {
					break
				}
				k := string(bytes.Trim(kv.Key[:], "\x00"))
				value := string(bytes.Trim(kv.Value[:], "\x00"))
				if k == key {
					return &value
				}
			}
		}
	}

	return generic.FindInFiles(lt.dataFiles.Files(), key)
}

func (lt *LSMTree) Clean() error {
	err := os.RemoveAll(memTableLogFilePath)
	if err != nil {
		return err
	}
	lt.memTable = mem_table.NewMemTable(prepareDataFile())
	return nil
}

func (lt *LSMTree) Index() error {
	kv := kv_protocol.KV{}
	lt.indexes = map[string]map[string]bytesRange{}

	for _, file := range lt.dataFiles.Files() {
		currentSymbol := ""
		startOffset, lengthOffset := 0, 0
		_, _ = file.Seek(0, io.SeekStart)
		i := 0
		lt.indexes[file.Name()] = map[string]bytesRange{}
		for {
			err := binary.Read(file, binary.LittleEndian, &kv)
			if err == io.EOF {
				break
			}
			key := string(bytes.Trim(kv.Key[:], "\x00"))
			firstLetter := string([]rune(key)[0])
			if firstLetter != currentSymbol {
				if i == 0 {
					startOffset = lengthOffset
				} else {
					lt.indexes[file.Name()][currentSymbol] = bytesRange{
						minByte: startOffset,
						maxByte: lengthOffset,
					}
					startOffset = lengthOffset
				}
				currentSymbol = firstLetter
			}
			lengthOffset += kv_protocol.KVLengthBytes
			i++
		}
		lt.indexes[file.Name()][currentSymbol] = bytesRange{
			minByte: startOffset,
			maxByte: lengthOffset,
		}
	}
	return nil
}
