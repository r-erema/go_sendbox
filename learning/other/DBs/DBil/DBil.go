package DBil

import (
	"bytes"
	"encoding/binary"
	"go_sendbox/learning/other/DBs/DBil/bloom_filter"
	"go_sendbox/learning/other/DBs/DBil/files_storage"
	"go_sendbox/learning/other/DBs/DBil/kv_protocol"
	"go_sendbox/learning/other/DBs/DBil/strategy"
	"go_sendbox/learning/other/DBs/DBil/strategy/hash_table"
	"go_sendbox/learning/other/DBs/DBil/strategy/lsm_tree"
	"io"
	"reflect"
)

const maxDataFileSizeBytes = 100000

var (
	bloomFilter  = bloom_filter.NewFilter()
	filesStorage = files_storage.NewDataFilesStorage()
)

type DBil struct {
	strategy     strategy.Strategy
	filesStorage *files_storage.DataFilesStorage
	bloomFilter  *bloom_filter.Filter
}

func NewLSMTreeBasedDB() *DBil {
	return &DBil{
		strategy:     lsm_tree.NewLSMTree(filesStorage, 5),
		filesStorage: filesStorage,
		bloomFilter:  bloomFilter,
	}
}

func NewHashTableDB() *DBil {
	return &DBil{
		strategy:     hash_table.NewHashTable(filesStorage),
		filesStorage: filesStorage,
		bloomFilter:  bloomFilter,
	}
}

func (db *DBil) Save(key, value string) error {
	err := db.strategy.Save(key, value)
	if err != nil {
		return err
	}
	db.bloomFilter.Add(key)

	file := db.filesStorage.LastFile()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() > maxDataFileSizeBytes {
		db.filesStorage.AppendNewFile()
	}

	return nil
}

func (db *DBil) Get(key string) *string {
	if !db.bloomFilter.Get(key) {
		return nil
	}

	return db.strategy.Get(key)
}

func (db *DBil) Clean() error {
	db.bloomFilter = bloom_filter.NewFilter()

	err := db.strategy.Clean()
	if err != nil {
		return err
	}

	err = db.filesStorage.Reset()
	if err != nil {
		return err
	}
	return nil
}

func (db *DBil) Index() error {
	return db.strategy.Index()
}

func (db *DBil) Compaction() error {
	vales := make(map[string]string)
	kv := kv_protocol.KV{}
	for _, file := range db.filesStorage.Files() {
		_, _ = file.Seek(0, io.SeekStart)
		for {
			err := binary.Read(file, binary.LittleEndian, &kv)
			if err == io.EOF {
				break
			}
			keyFromFile := string(bytes.Trim(kv.Key[:], "\x00"))
			valueFromFile := string(bytes.Trim(kv.Value[:], "\x00"))
			vales[keyFromFile] = valueFromFile
		}
	}

	_ = db.Clean()
	for key, value := range vales {
		err := db.Save(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DBil) MergeFiles() {
	tmpBuffer := map[string]string{}
	kv := kv_protocol.KV{}
	for _, file := range db.filesStorage.Files() {
		_, _ = file.Seek(0, io.SeekStart)
		for {
			err := binary.Read(file, binary.LittleEndian, &kv)
			if err == io.EOF {
				break
			}
			key := string(bytes.Trim(kv.Key[:], "\x00"))
			value := string(bytes.Trim(kv.Value[:], "\x00"))
			tmpBuffer[key] = value
		}
	}

	_ = db.Clean()

	switch reflect.TypeOf(db.strategy).String() {
	case "*hash_table.hashTable":
		db.strategy = hash_table.NewHashTable(db.filesStorage)
	case "*lsm_tree.LSMTree":
		db.strategy = lsm_tree.NewLSMTree(db.filesStorage, 5)
	}

	for key, value := range tmpBuffer {
		copy(kv.Key[:], key)
		copy(kv.Value[:], value)
		lastFile := db.filesStorage.LastFile()
		_ = binary.Write(lastFile, binary.LittleEndian, kv)
		db.bloomFilter.Add(key)
	}

}
