package cache

import (
	"os"
	"path/filepath"
	"time"
)

type (
	// Cache stores data of any type as JSON files, one file per key.
	Cache struct {
		maxAge int64
		dir    string
	}

	MetaData struct {
		Updated    time.Time
		NextUpdate time.Time
	}

	record[T any] struct {
		Meta MetaData
		Data T
	}
)

func (c *Cache) path(key string) string {
	return filepath.Join(c.dir, key+".json")
}

func newMetaData(maxAge int64) MetaData {
	return MetaData{
		Updated:    time.Now(),
		NextUpdate: time.Now().Add(time.Duration(maxAge) * time.Minute),
	}
}

// NewCache creates a cache layer instance backed by JSON files in dir.
// maxAge is the time in minutes until data needs a refresh.
func NewCache(maxAge int, dir string) (*Cache, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return &Cache{
		maxAge: int64(maxAge),
		dir:    dir,
	}, nil
}
