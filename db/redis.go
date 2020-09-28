package db

import (
	"time"
)

// StoreCache stores key/value to the storage with expiration time
func (s *datastore) StoreCache(key string, payload interface{}, exp time.Duration) error {
	err := s.rdb.Set(ctx, key, payload, exp).Err()
	if err != nil {
		return err
	}

	return nil
}

// DeleteCache removes key from the storage
func (s *datastore) DeleteCache(key string) (int64, error) {
	return s.rdb.Del(ctx, key).Result()
}
