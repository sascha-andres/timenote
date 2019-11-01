package cache

import (
	"fmt"
	"go.etcd.io/bbolt"
	"path"
	"time"
)

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

func NewCache(maxAge int, pathToDatabase string) (*Cache, error) {
	db, err := bbolt.Open(path.Join(pathToDatabase, "db"), 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return &Cache{
		maxAge: int64(maxAge),

		db: db,
	}, nil
}

func setupDatabase(db *bbolt.DB) error {
	if err := createBucket(db, "projects"); err != nil {
		return err
	}
	if err := createBucket(db, "clients"); err != nil {
		return err
	}
	return nil
}

func createBucket(db *bbolt.DB, name string) error {
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}
