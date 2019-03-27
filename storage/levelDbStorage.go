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
func (s *LevelDbStorage) Get(key string) []byte {
	var value []byte

	value, err := s.Db.Get([]byte(key), nil)
	if err != nil {
		logger.Error("db get data error", zap.Error(err))
	}
	return value
}

// GetAll will return all key-value in DB.
func (s *LevelDbStorage) GetAll() []string {
	var result []string

	iter := s.Db.NewIterator(nil, nil)
	for iter.Next() {
		result = append(result, string(iter.Key()))
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

// Close will close the DB.
func (s *LevelDbStorage) Close() error {
	err := s.Db.Close()
	if err != nil {
		logger.Error("db close error", zap.Error(err))
	}
	return err
}
