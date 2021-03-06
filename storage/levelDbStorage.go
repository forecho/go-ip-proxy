package storage

import (
	"errors"
	"github.com/syndtr/goleveldb/leveldb"
	"go-ip-proxy/logger"
	"go.uber.org/zap"
)

type LevelDbStorage struct {
	Db *leveldb.DB
}

// NewLevelDbStorage will return a leveldb object and error.
func NewLevelDbStorage(fileName string) (*LevelDbStorage, error) {
	if fileName == "" {
		return nil, errors.New("open leveldb whose fileName is empty")
	}
	db, err := leveldb.OpenFile(fileName, nil)
	if err != nil {
		return nil, err
	}

	storage := &LevelDbStorage{
		Db: db,
	}

	return storage, nil
}

// Get will get the json byte value of key.
func (s *LevelDbStorage) Get() string {
	var value string

	iter := s.Db.NewIterator(nil, nil)

	if iter.First() {
		return string(iter.Key())
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		logger.Error("db get one data error", zap.Error(err))
	}
	return value
}

// GetAll will return all key-value in DB.
func (s *LevelDbStorage) GetAll() []string {
	var result []string

	iter := s.Db.NewIterator(nil, nil)
	for iter.Next() {
		key := string(iter.Key())
		result = append(result, key)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		logger.Error("db get all data error", zap.Error(err))
	}
	return result
}

func (s *LevelDbStorage) Create(key string, value string) error {
	err := s.Db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		logger.Error("db create error", zap.Error(err))
	}
	return err
}

// Delete the value by the given key.
func (s *LevelDbStorage) Delete(key string) bool {
	isSucceed := false
	err := s.Db.Delete([]byte(key), nil)

	if err == nil {
		isSucceed = true
	}

	return isSucceed
}

// Close will close the DB.
func (s *LevelDbStorage) Close() {
	err := s.Db.Close()
	if err != nil {
		logger.Error("db close error", zap.Error(err))
	}
}
