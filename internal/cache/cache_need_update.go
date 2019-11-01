package cache

import (
	"fmt"
	"go.etcd.io/bbolt"
	"gopkg.in/yaml.v2"
	"time"
)

// NeedUpdate returns true if cache needs a refresh
func (c *Cache) NeedUpdate(workspace int) bool {
	return c.needUpdate(workspace, "projects") || c.needUpdate(workspace, "clients")
}

func (c *Cache) needUpdate(workspace int, bucket string) (needUpdate bool) {
	_ = c.db.Update(func(tx *bbolt.Tx) error {
		bucketName := fmt.Sprintf("%6d-%s", workspace, bucket)
		fmt.Println(bucketName)
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		v := bucket.Get([]byte("_meta"))
		needUpdate = len(v) == 0 || checkIfUpdateRequired(v)
		return nil
	})
	return
}

func checkIfUpdateRequired(data []byte) bool {
	var m MetaData
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		return true
	}
	return time.Now().Sub(m.NextUpdate) > 0
}
