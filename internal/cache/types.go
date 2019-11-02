package cache

import (
	"go.etcd.io/bbolt"
	"path"
	"time"
)

const clientBucketNameTemplate = "%6d-clients"
const projectBucketNameTemplate = "%6d-projects"
const metaKeyName = "_meta"
const allKeyName = "_all"

type (
	Cache struct {
		maxAge int64

		db *bbolt.DB
	}

	MetaData struct {
		Updated    time.Time
		NextUpdate time.Time
	}
)

// NewCache creates a cache layer instance and returns it
// inner workings like options are not yet available
func NewCache(maxAge int, pathToDatabase string) (*Cache, error) {
	db, err := bbolt.Open(path.Join(pathToDatabase, "db"), 0600, &bbolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, err
	}

	return &Cache{
		maxAge: int64(maxAge),

		db: db,
	}, nil
}
